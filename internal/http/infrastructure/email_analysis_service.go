// Package infrastructure はHTTP通信のインフラストラクチャ層を提供します。
// このファイルはOpenAI APIを使用したメール分析サービスを実装します。
package infrastructure

import (
	"business/internal/openai/domain"
	openair "business/internal/openai/infrastructure"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OpenAIRequest はOpenAI APIリクエストの構造体です
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message はOpenAI APIのメッセージ構造体です
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse はOpenAI APIレスポンスの構造体です
type OpenAIResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice はOpenAI APIの選択肢構造体です
type Choice struct {
	Message Message `json:"message"`
}

// APIError はOpenAI APIエラーの構造体です
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// EmailAnalysisServiceImpl はメール分析サービスの実装です
type EmailAnalysisServiceImpl struct {
	promptService openair.PromptService
}

// NewEmailAnalysisService はメール分析サービスを作成します
func NewEmailAnalysisService(promptService openair.PromptService) *EmailAnalysisServiceImpl {
	return &EmailAnalysisServiceImpl{
		promptService: promptService,
	}
}

// AnalyzeEmail はメールを分析します
func (s *EmailAnalysisServiceImpl) AnalyzeEmail(ctx context.Context, request *domain.EmailAnalysisRequest) (*domain.EmailAnalysisResult, error) {
	// リクエストの妥当性チェック
	if err := request.IsValid(); err != nil {
		return nil, fmt.Errorf("リクエスト妥当性エラー: %w", err)
	}

	// プロンプトファイルの読み込み
	promptText, err := s.promptService.LoadPrompt("text_analysis_prompt.txt")
	if err != nil {
		return nil, fmt.Errorf("プロンプト読み込みエラー: %w", err)
	}

	// プロンプトとメール本文を結合
	combinedText := promptText + "\n\n" + request.EmailText

	// 直接HTTPリクエストを送信
	content, err := s.sendOpenAIRequest(ctx, combinedText)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API呼び出しエラー: %w", err)
	}

	// JSONをパース
	var analysisData struct {
		MailCategory        string   `json:"メール区分"`
		StartPeriod         []string `json:"開始時期"`
		EndPeriod           string   `json:"終了時期"`
		WorkLocation        string   `json:"勤務場所"`
		PriceFrom           *int     `json:"単価FROM"`
		PriceTo             *int     `json:"単価TO"`
		Languages           []string `json:"言語"`
		Frameworks          []string `json:"フレームワーク"`
		RequiredSkillsMust  []string `json:"求めるスキル MUST"`
		RequiredSkillsWant  []string `json:"求めるスキル WANT"`
		RemoteWorkCategory  string   `json:"リモートワーク区分"`
		RemoteWorkFrequency *string  `json:"リモートワークの頻度"`
	}

	if err := json.Unmarshal([]byte(content), &analysisData); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %w, content: %s", err, content)
	}

	// 結果を作成
	result := &domain.EmailAnalysisResult{
		Subject:             request.Subject,
		MailCategory:        analysisData.MailCategory,
		StartPeriod:         analysisData.StartPeriod,
		EndPeriod:           analysisData.EndPeriod,
		WorkLocation:        analysisData.WorkLocation,
		PriceFrom:           analysisData.PriceFrom,
		PriceTo:             analysisData.PriceTo,
		Languages:           analysisData.Languages,
		Frameworks:          analysisData.Frameworks,
		RequiredSkillsMust:  analysisData.RequiredSkillsMust,
		RequiredSkillsWant:  analysisData.RequiredSkillsWant,
		RemoteWorkCategory:  analysisData.RemoteWorkCategory,
		RemoteWorkFrequency: analysisData.RemoteWorkFrequency,
	}

	return result, nil
}

// sendOpenAIRequest はOpenAI APIにリクエストを送信します
func (s *EmailAnalysisServiceImpl) sendOpenAIRequest(ctx context.Context, content string) (string, error) {
	// OpenAI APIリクエストを作成
	openaiRequest := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
		Temperature: 0.1,
		MaxTokens:   1000,
	}

	// リクエストボディをJSONに変換
	requestBody, err := json.Marshal(openaiRequest)
	if err != nil {
		return "", fmt.Errorf("リクエストJSONエンコードエラー: %w", err)
	}

	// HTTPリクエストを作成
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト作成エラー: %w", err)
	}

	// ヘッダーを設定
	httpRequest.Header.Set("Content-Type", "application/json")
	// OpenAI APIキーは環境変数から取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY環境変数が設定されていません")
	}
	httpRequest.Header.Set("Authorization", "Bearer "+apiKey)

	// HTTPクライアントを作成
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// リクエストを送信
	response, err := client.Do(httpRequest)
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト送信エラー: %w", err)
	}
	defer response.Body.Close()

	// レスポンスボディを読み取り
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み取りエラー: %w", err)
	}

	// HTTPステータスコードをチェック
	if response.StatusCode != http.StatusOK {
		var apiError OpenAIResponse
		if err := json.Unmarshal(responseBody, &apiError); err == nil && apiError.Error != nil {
			return "", fmt.Errorf("OpenAI APIエラー: %s", apiError.Error.Message)
		}
		return "", fmt.Errorf("OpenAI API HTTPエラー: %d", response.StatusCode)
	}

	// OpenAIレスポンスを解析
	var openaiResponse OpenAIResponse
	if err := json.Unmarshal(responseBody, &openaiResponse); err != nil {
		return "", fmt.Errorf("OpenAIレスポンスJSONデコードエラー: %w", err)
	}

	if len(openaiResponse.Choices) == 0 {
		return "", fmt.Errorf("OpenAIレスポンスに選択肢がありません")
	}

	return openaiResponse.Choices[0].Message.Content, nil
}
