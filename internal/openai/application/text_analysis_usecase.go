// Package application はテキスト字句解析のアプリケーション層を提供します。
// このファイルはメール本文のAI字句解析に関するユースケースを実装します。
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	ed "business/internal/emailstore/domain"
	gd "business/internal/gmail/domain"
	"business/internal/openai/domain"
	r "business/internal/openai/infrastructure"
)

// TextAnalysisUseCaseImpl はテキスト字句解析のユースケース実装です
type TextAnalysisUseCaseImpl struct {
	textAnalysisService TextAnalysisService
	promptService       r.PromptService
}

// NewTextAnalysisUseCase はテキスト字句解析ユースケースを作成します
func NewTextAnalysisUseCase(textAnalysisService TextAnalysisService, promptService r.PromptService) TextAnalysisUseCase {
	return &TextAnalysisUseCaseImpl{
		textAnalysisService: textAnalysisService,
		promptService:       promptService,
	}
}

// AnalyzeText はテキストをAIで字句解析します
func (u *TextAnalysisUseCaseImpl) AnalyzeText(ctx context.Context, text string) (*domain.TextAnalysisResult, error) {
	// 空のテキストチェック
	if text == "" {
		return nil, domain.ErrEmptyText
	}

	// テキスト解析リクエストを作成
	request := domain.NewTextAnalysisRequest(text)

	// AI解析実行
	result, err := u.textAnalysisService.AnalyzeText(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("テキスト解析エラー: %w", err)
	}

	return result, nil
}

// AnalyzeTextWithOptions はオプション付きでテキストをAIで字句解析します
func (u *TextAnalysisUseCaseImpl) AnalyzeTextWithOptions(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error) {
	// リクエストの妥当性チェック
	if err := request.IsValid(); err != nil {
		return nil, fmt.Errorf("リクエスト妥当性エラー: %w", err)
	}

	// AI解析実行
	result, err := u.textAnalysisService.AnalyzeText(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("テキスト解析エラー: %w", err)
	}

	return result, nil
}

// AnalyzeEmailText はメール本文をAIで字句解析します（後方互換性のため残す）
// 1. プロンプトファイルの内容を読み込み
// 2. プロンプト + メール本文を結合
// 3. AI APIに送信して解析結果を取得
func (u *TextAnalysisUseCaseImpl) AnalyzeEmailText(ctx context.Context, email gd.GmailMessage) ([]ed.AnalysisResult, error) {
	// 空のメール本文チェック
	emailText := email.Body
	messageID := email.ID
	subject := email.Subject
	if emailText == "" {
		return nil, domain.ErrEmptyText
	}

	// プロンプトファイルの読み込み
	promptText, err := u.promptService.LoadPrompt("text_analysis_prompt.txt")
	if err != nil {
		return nil, fmt.Errorf("プロンプト読み込みエラー: %w", err)
	}

	// プロンプトとメール本文を結合
	combinedText := promptText + "\n\n" + emailText

	// テキスト解析リクエストを作成
	request := domain.NewTextAnalysisRequest(combinedText)
	request.Metadata["source"] = "email"
	request.Metadata["message_id"] = messageID
	request.Metadata["subject"] = subject

	// リクエストの妥当性チェック
	if err := request.IsValid(); err != nil {
		return nil, fmt.Errorf("リクエスト妥当性エラー: %w", err)
	}

	// AI解析実行
	response, err := u.textAnalysisService.AnalyzeTextMultiple(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("テキスト解析エラー: %w", err)
	}

	var results []ed.AnalysisResult
	// 結果にメタデータを設定
	for _, analysisResult := range response {
		analysisResult.MessageID = messageID
		analysisResult.Subject = subject
		analysisResult.AnalyzedAt = time.Now()

		project, err := u.convertToProjectAnalysisResults(email, analysisResult)
		if err != nil {
			log.Printf("案件変換失敗: %v", err)
			continue
		}
		results = append(results, project...)
	}

	return results, nil
}

// convertToProjectAnalysisResults はTextAnalysisResultを案件分析結果に変換します
func (u *TextAnalysisUseCaseImpl) convertToProjectAnalysisResults(email gd.GmailMessage, result *domain.TextAnalysisResult) ([]ed.AnalysisResult, error) {
	// openai_responseキーから配列形式のレスポンスを処理
	openaiResponse, exists := result.RawResponse["openai_response"]
	if !exists {
		fmt.Printf("案件変換失敗 (GmailID: %s): \n", email.ID)
		return []ed.AnalysisResult{}, fmt.Errorf("OpenAIのレスポンスに 'openai_response' が存在しません。")
	}

	// レスポンスの型確認
	responseStr, ok := openaiResponse.(string)
	if !ok {
		fmt.Printf("案件変換失敗 (GmailID: %s): \n", email.ID)
		return []ed.AnalysisResult{}, fmt.Errorf("openai_responseが文字列ではありません。")
	}

	// JSON配列をパース
	var analysisResults []ed.AnalysisResult
	if err := json.Unmarshal([]byte(responseStr), &analysisResults); err != nil {
		fmt.Printf("案件変換失敗 (GmailID: %s): \n", email.ID)
		return []ed.AnalysisResult{}, fmt.Errorf("JSON配列のパースに失敗しました: %w", err)
	}

	if len(analysisResults) == 0 {
		return []ed.AnalysisResult{}, fmt.Errorf("案件情報が見つかりません。(GmailID: %s)", email.ID)
	}

	var projects []ed.AnalysisResult
	for _, analysisResult := range analysisResults {

		fmt.Printf("StartDates toStringSlice: %v \n", strings.Join(analysisResult.StartPeriod, ", "))

		project := ed.AnalysisResult{
			ProjectName:         result.Summary,
			GmailID:             email.ID,
			Summary:             analysisResult.ProjectName,
			Subject:             email.Subject,
			From:                email.ExtractSenderName(),
			FromEmail:           email.ExtractEmailAddress(),
			Body:                email.Body,
			ReceivedDate:        email.Date,
			Category:            analysisResult.Category,
			StartPeriod:         analysisResult.StartPeriod,
			EndPeriod:           analysisResult.EndPeriod,
			WorkLocation:        analysisResult.WorkLocation,
			PriceFrom:           analysisResult.PriceFrom,
			PriceTo:             analysisResult.PriceTo,
			Languages:           analysisResult.Languages,
			Frameworks:          analysisResult.Frameworks,
			Positions:           analysisResult.Positions,
			WorkTypes:           analysisResult.WorkTypes,
			RequiredSkillsMust:  analysisResult.RequiredSkillsMust,
			RequiredSkillsWant:  analysisResult.RequiredSkillsWant,
			RemoteWorkCategory:  analysisResult.RemoteWorkCategory,
			RemoteWorkFrequency: analysisResult.RemoteWorkFrequency,
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// DisplayAnalysisResult は解析結果をターミナルに表示します
func (u *TextAnalysisUseCaseImpl) DisplayAnalysisResult(result *domain.TextAnalysisResult) error {
	if result == nil {
		return fmt.Errorf("解析結果がnilです")
	}

	fmt.Println("=== テキスト字句解析結果 ===")
	fmt.Printf("メッセージID: %s\n", result.MessageID)
	fmt.Printf("件名: %s\n", result.Subject)
	fmt.Printf("解析日時: %s\n", result.AnalyzedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("言語: %s\n", result.Language)
	fmt.Printf("信頼度: %.2f\n", result.Confidence)
	fmt.Println()

	// 感情分析結果
	fmt.Println("--- 感情分析 ---")
	fmt.Printf("スコア: %.2f\n", result.Sentiment.Score)
	fmt.Printf("ラベル: %s\n", result.Sentiment.Label)
	fmt.Printf("信頼度: %.2f\n", result.Sentiment.Confidence)
	fmt.Println()

	// キーワード
	if len(result.Keywords) > 0 {
		fmt.Println("--- キーワード ---")
		for i, keyword := range result.Keywords {
			fmt.Printf("%d. %s (関連度: %.2f, 出現回数: %d)\n",
				i+1, keyword.Text, keyword.Relevance, keyword.Count)
		}
		fmt.Println()
	}

	// エンティティ
	if len(result.Entities) > 0 {
		fmt.Println("--- エンティティ ---")
		for i, entity := range result.Entities {
			fmt.Printf("%d. %s (%s) - 重要度: %.2f\n",
				i+1, entity.Name, entity.Type, entity.Salience)
		}
		fmt.Println()
	}

	// 要約
	if result.Summary != "" {
		fmt.Println("--- 要約 ---")
		fmt.Println(result.Summary)
		fmt.Println()
	}

	// カテゴリ
	if len(result.Categories) > 0 {
		fmt.Println("--- カテゴリ ---")
		for i, category := range result.Categories {
			fmt.Printf("%d. %s (信頼度: %.2f)\n",
				i+1, category.Name, category.Confidence)
		}
		fmt.Println()
	}

	return nil
}

func toStringSlice(input interface{}) []string {
	result := []string{}
	if arr, ok := input.([]interface{}); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
	}
	return result
}

func toIntPtr(input interface{}) *int {
	switch v := input.(type) {
	case float64:
		if v > 0 {
			i := int(v)
			return &i
		}
	case int:
		if v > 0 {
			return &v
		}
	}
	return nil
}
