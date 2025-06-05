// Package application はメール分析のアプリケーション層を提供します。
// このファイルはメール分析に関するユースケースを実装します。
package application

import (
	authdomain "business/internal/gmail/domain"
	"business/internal/openai/domain"
	r "business/internal/openai/infrastructure"
	"context"
	"fmt"
	"strings"
)

// EmailAnalysisUseCaseImpl はメール分析のユースケース実装です
type EmailAnalysisUseCaseImpl struct {
	emailAnalysisService r.EmailAnalysisService
}

// NewEmailAnalysisUseCase はメール分析ユースケースを作成します
func NewEmailAnalysisUseCase(emailAnalysisService r.EmailAnalysisService) EmailAnalysisUseCase {
	return &EmailAnalysisUseCaseImpl{
		emailAnalysisService: emailAnalysisService,
	}
}

// AnalyzeEmailContent はメール内容を分析します
func (u *EmailAnalysisUseCaseImpl) AnalyzeEmailContent(ctx context.Context, message *authdomain.GmailMessage) (*domain.EmailAnalysisResult, error) {
	// 空のメール本文チェック
	if message.Body == "" {
		return nil, domain.ErrEmptyEmailText
	}

	// メール分析リクエストを作成
	request := domain.NewEmailAnalysisRequest(message.Body, message.ID, message.Subject)

	// AI分析実行
	result, err := u.emailAnalysisService.AnalyzeEmail(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("メール分析エラー: %w", err)
	}

	// Gmailの情報を設定
	result.GmailID = message.ID
	result.Subject = message.Subject
	result.From = message.From
	result.FromEmail = extractEmailAddress(message.From)
	result.Date = message.Date
	result.Body = message.Body

	return result, nil
}

// DisplayEmailAnalysisResult はメール分析結果をターミナルに表示します
func (u *EmailAnalysisUseCaseImpl) DisplayEmailAnalysisResult(result *domain.EmailAnalysisResult) error {
	if result == nil {
		return fmt.Errorf("分析結果がnilです")
	}

	fmt.Println("=== メール分析結果 ===")
	fmt.Printf("\"ID\": \"%s\"\n", result.GmailID)
	fmt.Printf("\"件名\": \"%s\"\n", result.Subject)
	fmt.Printf("\"差出人名\": \"%s\"\n", result.From)
	fmt.Printf("\"メールアドレス\": \"%s\"\n", result.FromEmail)
	fmt.Printf("\"日時\": \"%s\"\n", result.Date.Format("2006/01/02"))
	// fmt.Printf("\"本文\": \"%s\"\n", truncateText(result.Body, 200))
	fmt.Printf("\"メール区分\": \"%s\"\n", result.MailCategory)

	// 配列の表示
	fmt.Printf("\"開始時期\": %s\n", formatStringArray(result.StartPeriod))
	fmt.Printf("\"終了時期\": \"%s\"\n", result.EndPeriod)
	fmt.Printf("\"勤務場所\": \"%s\"\n", result.WorkLocation)

	// 数値の表示（nullの場合は"null"と表示）
	if result.PriceFrom != nil {
		fmt.Printf("\"単価FROM\": %d\n", *result.PriceFrom)
	} else {
		fmt.Printf("\"単価FROM\": null\n")
	}

	if result.PriceTo != nil {
		fmt.Printf("\"単価TO\": %d\n", *result.PriceTo)
	} else {
		fmt.Printf("\"単価TO\": null\n")
	}

	fmt.Printf("\"言語\": %s\n", formatStringArray(result.Languages))
	fmt.Printf("\"フレームワーク\": %s\n", formatStringArray(result.Frameworks))
	fmt.Printf("\"求めるスキル MUST\": %s\n", formatStringArray(result.RequiredSkillsMust))
	fmt.Printf("\"求めるスキル WANT\": %s\n", formatStringArray(result.RequiredSkillsWant))
	fmt.Printf("\"リモートワーク区分\": \"%s\"\n", result.RemoteWorkCategory)

	if result.RemoteWorkFrequency != nil {
		fmt.Printf("\"リモートワークの頻度\": \"%s\"\n", *result.RemoteWorkFrequency)
	} else {
		fmt.Printf("\"リモートワークの頻度\": null\n")
	}

	return nil
}

// extractEmailAddress は差出人文字列からメールアドレスを抽出します
func extractEmailAddress(from string) string {
	// "Name <email@example.com>" 形式からメールアドレスを抽出
	if strings.Contains(from, "<") && strings.Contains(from, ">") {
		start := strings.Index(from, "<")
		end := strings.Index(from, ">")
		if start < end {
			email := from[start+1 : end]
			// 抽出したメールアドレスが有効かチェック
			if strings.Contains(email, "@") {
				return email
			}
		}
	}

	// メールアドレスのみの場合はそのまま返す
	if strings.Contains(from, "@") {
		return from
	}

	return ""
}

// truncateText はテキストを指定された長さで切り詰めます
func truncateText(text string, maxLen int) string {
	// ルーン（文字）単位で処理して日本語文字を正しく扱う
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen]) + "..."
}

// formatStringArray は文字列配列をJSON形式で表示します
func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}

	var quoted []string
	for _, s := range arr {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", s))
	}

	return "[" + strings.Join(quoted, ", ") + "]"
}
