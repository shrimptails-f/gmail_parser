// Package application は認証機能のアプリケーション層を提供します。
// このファイルはGmailメッセージ取得ユースケースの実装を定義します。
package application

import (
	"business/internal/gmail/domain"
	r "business/internal/gmail/infrastructure"
	"context"
	"errors"
	"fmt"
	"time"
)

// GmailMessageUseCaseImpl はGmailMessageUseCaseの実装です
type GmailMessageUseCaseImpl struct {
	gmailAuthService    r.GmailAuthService
	gmailMessageService r.GmailMessageService
}

// NewGmailMessageUseCase はGmailMessageUseCaseの新しいインスタンスを作成します
func NewGmailMessageUseCase(gmailAuthService r.GmailAuthService, gmailMessageService r.GmailMessageService) GmailMessageUseCase {
	return &GmailMessageUseCaseImpl{
		gmailAuthService:    gmailAuthService,
		gmailMessageService: gmailMessageService,
	}
}

// GetMessages はメッセージ一覧を取得します
func (u *GmailMessageUseCaseImpl) GetMessages(ctx context.Context, config domain.GmailAuthConfig, maxResults int64) ([]domain.GmailMessage, error) {
	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 認証情報を取得
	credential, err := u.gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err != nil {
		return nil, fmt.Errorf("認証情報の読み込みに失敗しました: %w", err)
	}

	// 認証情報の有効性をチェック
	if !credential.IsValid() {
		return nil, errors.New("認証情報が無効です。再認証が必要です")
	}

	// メッセージ一覧を取得
	messages, err := u.gmailMessageService.GetMessages(ctx, *credential, config.ApplicationName, maxResults)
	if err != nil {
		return nil, fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	return messages, nil
}

// GetAllMessagesByLabelPathFromToday は指定されたラベルパスの当日0時以降のメッセージを全件取得します
func (u *GmailMessageUseCaseImpl) GetAllMessagesByLabelPathFromToday(ctx context.Context, config domain.GmailAuthConfig, labelPath string, maxResults int64) ([]domain.GmailMessage, error) {
	// ラベルパスの妥当性をチェック
	if labelPath == "" {
		return nil, errors.New("ラベルパスが指定されていません")
	}

	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 認証情報を取得
	credential, err := u.gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err != nil {
		return nil, fmt.Errorf("認証情報の読み込みに失敗しました: %w", err)
	}

	// 認証情報の有効性をチェック
	if !credential.IsValid() {
		return nil, errors.New("認証情報が無効です。再認証が必要です")
	}

	// ラベル一覧を取得
	labels, err := u.gmailMessageService.GetLabels(ctx, *credential, config.ApplicationName)
	if err != nil {
		return nil, fmt.Errorf("ラベル一覧の取得に失敗しました: %w", err)
	}

	// 指定されたラベルパスに一致するラベルIDを検索
	var labelID string
	for _, label := range labels {
		if label.Name == labelPath {
			labelID = label.ID
			break
		}
	}

	if labelID == "" {
		return nil, fmt.Errorf("指定されたラベル '%s' が見つかりません", labelPath)
	}

	// 当日の0時を取得
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 全メッセージを取得（ページネーションで0件になるまで取得）
	var allMessages []domain.GmailMessage
	pageToken := ""

	for {
		// 指定されたラベルと日付以降のメッセージを取得
		messages, nextPageToken, err := u.gmailMessageService.GetMessagesByLabelAndDate(ctx, *credential, config.ApplicationName, labelID, todayStart, maxResults, pageToken)
		if err != nil {
			return nil, fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
		}

		// 取得したメッセージを追加
		allMessages = append(allMessages, messages...)

		// 次のページがない、または取得件数が0の場合は終了
		if nextPageToken == "" || len(messages) == 0 {
			break
		}

		// 次のページトークンを設定
		pageToken = nextPageToken
	}

	return allMessages, nil
}

// GetMessage は指定されたIDのメッセージを取得します
func (u *GmailMessageUseCaseImpl) GetMessage(ctx context.Context, config domain.GmailAuthConfig, messageID string) (*domain.GmailMessage, error) {
	// メッセージIDの妥当性をチェック
	if messageID == "" {
		return nil, errors.New("メッセージIDが指定されていません")
	}

	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 認証情報を取得
	credential, err := u.gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err != nil {
		return nil, fmt.Errorf("認証情報の読み込みに失敗しました: %w", err)
	}

	// 認証情報の有効性をチェック
	if !credential.IsValid() {
		return nil, errors.New("認証情報が無効です。再認証が必要です")
	}

	// メッセージを取得
	message, err := u.gmailMessageService.GetMessage(ctx, *credential, config.ApplicationName, messageID)
	if err != nil {
		return nil, fmt.Errorf("メッセージの取得に失敗しました: %w", err)
	}

	return message, nil
}

// GetMessagesByLabelPath は指定されたラベルパスのメッセージ一覧を取得します
func (u *GmailMessageUseCaseImpl) GetMessagesByLabelPath(ctx context.Context, config domain.GmailAuthConfig, labelPath string, maxResults int64) ([]domain.GmailMessage, error) {
	// ラベルパスの妥当性をチェック
	if labelPath == "" {
		return nil, errors.New("ラベルパスが指定されていません")
	}

	// 設定の妥当性をチェック
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// 認証情報を取得
	credential, err := u.gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err != nil {
		return nil, fmt.Errorf("認証情報の読み込みに失敗しました: %w", err)
	}

	// 認証情報の有効性をチェック
	if !credential.IsValid() {
		return nil, errors.New("認証情報が無効です。再認証が必要です")
	}

	// ラベル一覧を取得
	labels, err := u.gmailMessageService.GetLabels(ctx, *credential, config.ApplicationName)
	if err != nil {
		return nil, fmt.Errorf("ラベル一覧の取得に失敗しました: %w", err)
	}

	// 指定されたラベルパスに一致するラベルIDを検索
	var labelID string
	for _, label := range labels {
		if label.Name == labelPath {
			labelID = label.ID
			break
		}
	}

	if labelID == "" {
		return nil, fmt.Errorf("指定されたラベル '%s' が見つかりません", labelPath)
	}

	// 指定されたラベルのメッセージ一覧を取得
	messages, err := u.gmailMessageService.GetMessagesByLabel(ctx, *credential, config.ApplicationName, labelID, maxResults)
	if err != nil {
		return nil, fmt.Errorf("ラベル指定メッセージ一覧の取得に失敗しました: %w", err)
	}

	return messages, nil
}
