// Package infrastructure はAI機能のインフラストラクチャ層を提供します。
// このファイルはOpenAI APIとの通信を行うサービスを実装します。
package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"business/internal/openai/domain"
)

// OpenAI APIの定数
const (
	OpenAIAPIURL   = "https://api.openai.com/v1/chat/completions"
	OpenAIModel    = "gpt-3.5-turbo"
	RequestTimeout = 30 * time.Second
	MaxTokens      = 2000
	Temperature    = 0.3
)

// OpenAIService はOpenAI APIを使用したテキスト解析サービスです
type OpenAIService struct {
	apiKey     string
	httpClient *http.Client
}

// NewOpenAIService はOpenAIサービスを作成します
func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

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

// AnalysisResponse はAI解析結果のレスポンス構造体です
type AnalysisResponse struct {
	Sentiment  SentimentResponse  `json:"sentiment"`
	Keywords   []KeywordResponse  `json:"keywords"`
	Entities   []EntityResponse   `json:"entities"`
	Summary    string             `json:"summary"`
	Categories []CategoryResponse `json:"categories"`
	Language   string             `json:"language"`
	Confidence float64            `json:"confidence"`
}

// SentimentResponse は感情分析レスポンスの構造体です
type SentimentResponse struct {
	Score      float64 `json:"score"`
	Magnitude  float64 `json:"magnitude"`
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

// KeywordResponse はキーワード抽出レスポンスの構造体です
type KeywordResponse struct {
	Text      string  `json:"text"`
	Relevance float64 `json:"relevance"`
	Count     int     `json:"count"`
	Category  string  `json:"category"`
}

// EntityResponse はエンティティ抽出レスポンスの構造体です
type EntityResponse struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Salience   float64 `json:"salience"`
	Confidence float64 `json:"confidence"`
}

// CategoryResponse はカテゴリ分類レスポンスの構造体です
type CategoryResponse struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	Path       string  `json:"path"`
}

// EmailProjectResponse は営業案件メール用のレスポンス構造体です
type EmailProjectResponse struct {
	EmailType       string   `json:"メール区分"`
	ProjectName     string   `json:"案件名"`
	Tasks           []string `json:"業務"`
	StartDates      []string `json:"入場時期・開始時期"`
	EndDate         string   `json:"終了時期"`
	WorkLocation    string   `json:"勤務場所"`
	SalaryFrom      int      `json:"単価FROM"`
	SalaryTo        int      `json:"単価TO"`
	Languages       []string `json:"言語"`
	Frameworks      []string `json:"フレームワーク"`
	Positions       []string `json:"ポジション"`
	RequiredSkills  []string `json:"求めるスキル MUST"`
	PreferredSkills []string `json:"求めるスキル WANT"`
	RemoteWorkType  string   `json:"リモートワーク区分"`
	RemoteFrequency string   `json:"リモートワークの頻度"`
}

// AnalyzeText はテキストをOpenAI APIで解析します
func (s *OpenAIService) AnalyzeText(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error) {
	// リクエストの妥当性チェック
	if err := request.IsValid(); err != nil {
		return nil, err
	}

	// プロンプトを構築
	prompt := s.buildPrompt(request)

	// OpenAI APIリクエストを作成
	openaiRequest := OpenAIRequest{
		Model: OpenAIModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   MaxTokens,
		Temperature: Temperature,
	}

	// APIリクエストを送信
	responseBody, err := s.sendRequest(ctx, openaiRequest)
	if err != nil {
		return nil, fmt.Errorf("OpenAI APIリクエストエラー: %w", err)
	}

	// レスポンスを解析
	result, err := s.parseResponse(responseBody)
	if err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	return result, nil
}

// parseResponseMultiple はOpenAIのレスポンスを複数の解析結果に変換します
func (s *OpenAIService) parseResponseMultiple(responseBody string) ([]*domain.TextAnalysisResult, error) {
	// まず営業案件メール用の配列として解析を試行
	var emailProjectArray []EmailProjectResponse
	if err := json.Unmarshal([]byte(responseBody), &emailProjectArray); err != nil {
		// 営業案件メール用の配列でない場合は、通常の解析レスポンス配列として試行
		var analysisResponseArray []AnalysisResponse
		if err := json.Unmarshal([]byte(responseBody), &analysisResponseArray); err != nil {
			// 配列でない場合は単一オブジェクトとして解析
			var singleResponse AnalysisResponse
			if err := json.Unmarshal([]byte(responseBody), &singleResponse); err != nil {
				return nil, fmt.Errorf("解析結果JSONデコードエラー: %w", err)
			}
			// 単一オブジェクトを配列に変換
			analysisResponseArray = []AnalysisResponse{singleResponse}
		}

		// 通常の解析レスポンスをドメインモデルに変換
		return s.convertAnalysisResponsesToResults(analysisResponseArray, responseBody), nil
	}

	// 営業案件メール用のレスポンスをドメインモデルに変換
	return s.convertEmailProjectsToResults(emailProjectArray, responseBody), nil
}

// convertAnalysisResponsesToResults は通常の解析レスポンスをドメインモデルに変換します
func (s *OpenAIService) convertAnalysisResponsesToResults(analysisResponseArray []AnalysisResponse, responseBody string) []*domain.TextAnalysisResult {
	var results []*domain.TextAnalysisResult
	for _, analysisResponse := range analysisResponseArray {
		result := &domain.TextAnalysisResult{
			AnalyzedAt: time.Now(),
			Sentiment: domain.SentimentAnalysis{
				Score:      analysisResponse.Sentiment.Score,
				Magnitude:  analysisResponse.Sentiment.Magnitude,
				Label:      analysisResponse.Sentiment.Label,
				Confidence: analysisResponse.Sentiment.Confidence,
			},
			Summary:    analysisResponse.Summary,
			Language:   analysisResponse.Language,
			Confidence: analysisResponse.Confidence,
			RawResponse: map[string]interface{}{
				"openai_response": responseBody,
			},
		}

		// キーワードを変換
		for _, kw := range analysisResponse.Keywords {
			result.Keywords = append(result.Keywords, domain.Keyword{
				Text:      kw.Text,
				Relevance: kw.Relevance,
				Count:     kw.Count,
				Category:  kw.Category,
			})
		}

		// エンティティを変換
		for _, entity := range analysisResponse.Entities {
			result.Entities = append(result.Entities, domain.Entity{
				Name:       entity.Name,
				Type:       entity.Type,
				Salience:   entity.Salience,
				Confidence: entity.Confidence,
			})
		}

		// カテゴリを変換
		for _, cat := range analysisResponse.Categories {
			result.Categories = append(result.Categories, domain.Category{
				Name:       cat.Name,
				Confidence: cat.Confidence,
				Path:       cat.Path,
			})
		}

		results = append(results, result)
	}
	return results
}

// convertEmailProjectsToResults は営業案件メール用のレスポンスをドメインモデルに変換します
func (s *OpenAIService) convertEmailProjectsToResults(emailProjects []EmailProjectResponse, responseBody string) []*domain.TextAnalysisResult {
	var results []*domain.TextAnalysisResult
	for _, project := range emailProjects {
		result := &domain.TextAnalysisResult{
			AnalyzedAt: time.Now(),
			Summary:    project.ProjectName,
			Language:   "ja",
			Confidence: 0.9, // 営業案件メールの場合は固定値
			RawResponse: map[string]interface{}{
				"openai_response": responseBody,
				"email_project":   project,
			},
		}

		// 営業案件メール特有の情報をキーワードとして追加
		if project.ProjectName != "" {
			result.Keywords = append(result.Keywords, domain.Keyword{
				Text:      project.ProjectName,
				Relevance: 1.0,
				Count:     1,
				Category:  "案件名",
			})
		}

		// 言語・フレームワークをキーワードとして追加
		for _, lang := range project.Languages {
			if lang != "" {
				result.Keywords = append(result.Keywords, domain.Keyword{
					Text:      lang,
					Relevance: 0.8,
					Count:     1,
					Category:  "言語",
				})
			}
		}

		for _, fw := range project.Frameworks {
			if fw != "" {
				result.Keywords = append(result.Keywords, domain.Keyword{
					Text:      fw,
					Relevance: 0.8,
					Count:     1,
					Category:  "フレームワーク",
				})
			}
		}

		// ポジションをエンティティとして追加
		for _, pos := range project.Positions {
			if pos != "" {
				result.Entities = append(result.Entities, domain.Entity{
					Name:       pos,
					Type:       "POSITION",
					Salience:   0.7,
					Confidence: 0.9,
				})
			}
		}

		results = append(results, result)
	}
	return results
}

// AnalyzeTextMultiple はテキストをOpenAI APIで解析し、複数の結果を返します
func (s *OpenAIService) AnalyzeTextMultiple(ctx context.Context, request *domain.TextAnalysisRequest) ([]*domain.TextAnalysisResult, error) {
	// リクエストの妥当性チェック
	if err := request.IsValid(); err != nil {
		return nil, err
	}

	// プロンプトを構築
	prompt := s.buildPrompt(request)

	// OpenAI APIリクエストを作成
	openaiRequest := OpenAIRequest{
		Model: OpenAIModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   MaxTokens,
		Temperature: Temperature,
	}

	// APIリクエストを送信
	responseBody, err := s.sendRequest(ctx, openaiRequest)
	if err != nil {
		return nil, fmt.Errorf("OpenAI APIリクエストエラー: %w", err)
	}

	// レスポンスを解析（複数案件対応）
	results, err := s.parseResponseMultiple(responseBody)
	if err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	return results, nil
}

// buildPrompt はリクエストからOpenAI用のプロンプトを構築します
// 複数案件対応のため、プロンプトファイルの内容をそのまま使用します
func (s *OpenAIService) buildPrompt(request *domain.TextAnalysisRequest) string {
	// メタデータからソースがemailの場合は、リクエストテキストをそのまま返す
	// （プロンプトファイルの内容は既にユースケース層で結合済み）
	if source, exists := request.Metadata["source"]; exists && source == "email" {
		return request.Text
	}

	// 通常のテキスト解析の場合は従来通りのプロンプトを構築
	var promptBuilder strings.Builder

	promptBuilder.WriteString("以下のテキストを字句解析して、JSON形式で結果を返してください。\n\n")
	promptBuilder.WriteString("テキスト:\n")
	promptBuilder.WriteString(request.Text)
	promptBuilder.WriteString("\n\n")

	promptBuilder.WriteString("以下の形式でJSON応答してください:\n")
	promptBuilder.WriteString("{\n")

	if request.Options.EnableSentiment {
		promptBuilder.WriteString(`  "sentiment": {
    "score": -1.0から1.0の感情スコア,
    "magnitude": 0.0以上の感情の強さ,
    "label": "POSITIVE", "NEGATIVE", "NEUTRAL", "MIXED"のいずれか,
    "confidence": 0.0から1.0の信頼度
  },` + "\n")
	}

	if request.Options.EnableKeywords {
		promptBuilder.WriteString(`  "keywords": [
    {
      "text": "キーワード",
      "relevance": 0.0から1.0の関連度,
      "count": 出現回数,
      "category": "カテゴリ"
    }
  ],` + "\n")
	}

	if request.Options.EnableEntities {
		promptBuilder.WriteString(`  "entities": [
    {
      "name": "エンティティ名",
      "type": "PERSON", "ORGANIZATION", "LOCATION"など,
      "salience": 0.0から1.0の重要度,
      "confidence": 0.0から1.0の信頼度
    }
  ],` + "\n")
	}

	if request.Options.EnableSummary {
		promptBuilder.WriteString(`  "summary": "テキストの要約",` + "\n")
	}

	if request.Options.EnableCategories {
		promptBuilder.WriteString(`  "categories": [
    {
      "name": "カテゴリ名",
      "confidence": 0.0から1.0の信頼度,
      "path": "カテゴリの階層パス"
    }
  ],` + "\n")
	}

	promptBuilder.WriteString(`  "language": "言語コード（ja, enなど）",
  "confidence": 0.0から1.0の全体的な信頼度
}

重要: 必ず有効なJSONのみを返してください。説明文は含めないでください。`)

	return promptBuilder.String()
}

// sendRequest はOpenAI APIにリクエストを送信します
func (s *OpenAIService) sendRequest(ctx context.Context, request OpenAIRequest) (string, error) {
	// リクエストボディをJSONに変換
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("リクエストJSONエンコードエラー: %w", err)
	}

	// HTTPリクエストを作成
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", OpenAIAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト作成エラー: %w", err)
	}

	// ヘッダーを設定
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+s.apiKey)

	// リクエストを送信
	response, err := s.httpClient.Do(httpRequest)
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

// parseResponse はOpenAIのレスポンスを解析結果に変換します
func (s *OpenAIService) parseResponse(responseBody string) (*domain.TextAnalysisResult, error) {
	var analysisResponse AnalysisResponse
	if err := json.Unmarshal([]byte(responseBody), &analysisResponse); err != nil {
		return nil, fmt.Errorf("解析結果JSONデコードエラー: %w", err)
	}

	// ドメインモデルに変換
	result := &domain.TextAnalysisResult{
		AnalyzedAt: time.Now(),
		Sentiment: domain.SentimentAnalysis{
			Score:      analysisResponse.Sentiment.Score,
			Magnitude:  analysisResponse.Sentiment.Magnitude,
			Label:      analysisResponse.Sentiment.Label,
			Confidence: analysisResponse.Sentiment.Confidence,
		},
		Summary:    analysisResponse.Summary,
		Language:   analysisResponse.Language,
		Confidence: analysisResponse.Confidence,
		RawResponse: map[string]interface{}{
			"openai_response": responseBody,
		},
	}

	// キーワードを変換
	for _, kw := range analysisResponse.Keywords {
		result.Keywords = append(result.Keywords, domain.Keyword{
			Text:      kw.Text,
			Relevance: kw.Relevance,
			Count:     kw.Count,
			Category:  kw.Category,
		})
	}

	// エンティティを変換
	for _, entity := range analysisResponse.Entities {
		result.Entities = append(result.Entities, domain.Entity{
			Name:       entity.Name,
			Type:       entity.Type,
			Salience:   entity.Salience,
			Confidence: entity.Confidence,
		})
	}

	// カテゴリを変換
	for _, cat := range analysisResponse.Categories {
		result.Categories = append(result.Categories, domain.Category{
			Name:       cat.Name,
			Confidence: cat.Confidence,
			Path:       cat.Path,
		})
	}

	return result, nil
}
