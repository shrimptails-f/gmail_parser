// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルはGmailメッセージ取得サービスの実装を定義します。
package infrastructure

import (
	"business/internal/gmail/domain"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// gmailMessageService はGmailメッセージ取得サービスの実装です
type gmailMessageService struct{}

// NewGmailMessageService は新しいGmailメッセージ取得サービスを作成します
func NewGmailMessageService() GmailMessageService {
	return &gmailMessageService{}
}

// GetLabels はラベル一覧を取得します
func (s *gmailMessageService) GetLabels(ctx context.Context, credential domain.GmailCredential, applicationName string) ([]domain.GmailInfo, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// ラベル一覧を取得
	response, err := service.Users.Labels.List("me").Do()
	if err != nil {
		return nil, fmt.Errorf("ラベル一覧の取得に失敗しました: %w", err)
	}

	// ドメインモデルに変換
	var labels []domain.GmailInfo
	for _, label := range response.Labels {
		labels = append(labels, domain.GmailInfo{
			ID:   label.Id,
			Name: label.Name,
			Type: label.Type,
		})
	}

	return labels, nil
}

// GetMessagesByLabel は指定されたラベルのメッセージ一覧を取得します
func (s *gmailMessageService) GetMessagesByLabel(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, maxResults int64) ([]domain.GmailMessage, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// 指定されたラベルのメッセージ一覧を取得
	call := service.Users.Messages.List("me").MaxResults(maxResults)
	if labelID != "" {
		call = call.LabelIds(labelID)
	}
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	// メッセージの詳細を取得
	var messages []domain.GmailMessage
	for _, msg := range response.Messages {
		message, err := s.GetMessage(ctx, credential, applicationName, msg.Id)
		if err != nil {
			// 個別のメッセージ取得に失敗した場合はログに記録して続行
			fmt.Printf("メッセージ %s の取得に失敗しました: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, *message)
	}

	return messages, nil
}

// GetMessages はメッセージ一覧を取得します
func (s *gmailMessageService) GetMessages(ctx context.Context, credential domain.GmailCredential, applicationName string, maxResults int64) ([]domain.GmailMessage, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// メッセージ一覧を取得
	call := service.Users.Messages.List("me").MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	// メッセージの詳細を取得
	var messages []domain.GmailMessage
	for _, msg := range response.Messages {
		message, err := s.GetMessage(ctx, credential, applicationName, msg.Id)
		if err != nil {
			// 個別のメッセージ取得に失敗した場合はログに記録して続行
			fmt.Printf("メッセージ %s の取得に失敗しました: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, *message)
	}

	return messages, nil
}

// GetMessage は指定されたIDのメッセージを取得します
func (s *gmailMessageService) GetMessage(ctx context.Context, credential domain.GmailCredential, applicationName string, messageID string) (*domain.GmailMessage, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// メッセージを取得
	message, err := service.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return nil, fmt.Errorf("メッセージの取得に失敗しました: %w", err)
	}

	// ドメインモデルに変換
	gmailMessage := s.convertToGmailMessage(message)
	return gmailMessage, nil
}

// createGmailService はGmail APIサービスを作成します
func (s *gmailMessageService) createGmailService(ctx context.Context, credential domain.GmailCredential, applicationName string) (*gmail.Service, error) {
	// OAuth2トークンを作成
	token := &oauth2.Token{
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		TokenType:    credential.TokenType,
		Expiry:       credential.ExpiresAt,
	}

	// OAuth2設定を作成（最小限の設定）
	config := &oauth2.Config{
		Endpoint: google.Endpoint,
	}

	// HTTPクライアントを作成
	client := config.Client(ctx, token)

	// Gmail APIサービスを作成
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	return service, nil
}

// convertToGmailMessage はGmail APIのメッセージをドメインモデルに変換します
func (s *gmailMessageService) convertToGmailMessage(message *gmail.Message) *domain.GmailMessage {
	gmailMessage := &domain.GmailMessage{
		ID: message.Id,
	}

	// ヘッダーから情報を抽出
	for _, header := range message.Payload.Headers {
		switch header.Name {
		case "Subject":
			gmailMessage.Subject = header.Value
		case "From":
			gmailMessage.From = header.Value
		case "To":
			gmailMessage.To = strings.Split(header.Value, ",")
			// トリム処理
			for i, to := range gmailMessage.To {
				gmailMessage.To[i] = strings.TrimSpace(to)
			}
		case "Date":
			// 日付をパース
			if date, err := time.Parse(time.RFC1123Z, header.Value); err == nil {
				gmailMessage.Date = date
			} else if date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", header.Value); err == nil {
				gmailMessage.Date = date
			} else {
				// パースに失敗した場合は現在時刻を設定
				gmailMessage.Date = time.Now()
			}
		}
	}

	// 本文を抽出
	body := s.extractBody(message.Payload)
	// HTMLタグを除去してプレーンテキストに変換
	gmailMessage.Body = stripHTMLTags(body)

	return gmailMessage
}

// extractBody はメッセージペイロードから本文を抽出します
func (s *gmailMessageService) extractBody(payload *gmail.MessagePart) string {
	// シンプルなテキスト本文の場合
	if payload.Body != nil && payload.Body.Data != "" {
		if decoded := s.decodeBase64Data(payload.Body.Data); decoded != "" {
			return decoded
		}
	}

	// マルチパートの場合
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			// text/plainを優先的に探す
			if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
				if decoded := s.decodeBase64Data(part.Body.Data); decoded != "" {
					return decoded
				}
			}
		}

		// text/plainが見つからない場合はtext/htmlを探す
		for _, part := range payload.Parts {
			if part.MimeType == "text/html" && part.Body != nil && part.Body.Data != "" {
				if decoded := s.decodeBase64Data(part.Body.Data); decoded != "" {
					return decoded
				}
			}
		}

		// 再帰的に探索
		for _, part := range payload.Parts {
			if body := s.extractBody(part); body != "" {
				return body
			}
		}
	}

	return ""
}

// GetMessagesByLabelWithPagination は指定されたラベルのメッセージ一覧をページネーションで取得します
func (s *gmailMessageService) GetMessagesByLabelWithPagination(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, maxResults int64, pageToken string) ([]domain.GmailMessage, string, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, "", fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// 指定されたラベルのメッセージ一覧を取得
	call := service.Users.Messages.List("me").MaxResults(maxResults)
	if labelID != "" {
		call = call.LabelIds(labelID)
	}
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	response, err := call.Do()
	if err != nil {
		return nil, "", fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	// メッセージの詳細を取得
	var messages []domain.GmailMessage
	for _, msg := range response.Messages {
		message, err := s.GetMessage(ctx, credential, applicationName, msg.Id)
		if err != nil {
			// 個別のメッセージ取得に失敗した場合はログに記録して続行
			fmt.Printf("メッセージ %s の取得に失敗しました: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, *message)
	}

	return messages, response.NextPageToken, nil
}

// GetMessagesByLabelAndDate は指定されたラベルと日付以降のメッセージ一覧を取得します
func (s *gmailMessageService) GetMessagesByLabelAndDate(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, afterDate time.Time, maxResults int64, pageToken string) ([]domain.GmailMessage, string, error) {
	// Gmail APIサービスを作成
	service, err := s.createGmailService(ctx, credential, applicationName)
	if err != nil {
		return nil, "", fmt.Errorf("Gmail APIサービスの作成に失敗しました: %w", err)
	}

	// 日付検索クエリを作成（YYYY/MM/DD形式）
	dateQuery := fmt.Sprintf("after:%s", afterDate.Format("2006/01/02"))

	// 指定されたラベルのメッセージ一覧を取得
	call := service.Users.Messages.List("me").MaxResults(maxResults).Q(dateQuery)
	if labelID != "" {
		call = call.LabelIds(labelID)
	}
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	response, err := call.Do()
	if err != nil {
		return nil, "", fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	// メッセージの詳細を取得
	var messages []domain.GmailMessage
	for _, msg := range response.Messages {
		message, err := s.GetMessage(ctx, credential, applicationName, msg.Id)
		if err != nil {
			// 個別のメッセージ取得に失敗した場合はログに記録して続行
			fmt.Printf("メッセージ %s の取得に失敗しました: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, *message)
	}

	return messages, response.NextPageToken, nil
}

// decodeBase64Data はbase64データをデコードします（複数の形式を試行）
func (s *gmailMessageService) decodeBase64Data(data string) string {
	// URL-safe base64を試行
	if decoded, err := base64.URLEncoding.DecodeString(data); err == nil {
		return string(decoded)
	}

	// 標準のbase64を試行
	if decoded, err := base64.StdEncoding.DecodeString(data); err == nil {
		return string(decoded)
	}

	// Raw URL-safe base64を試行
	if decoded, err := base64.RawURLEncoding.DecodeString(data); err == nil {
		return string(decoded)
	}

	return ""
}
