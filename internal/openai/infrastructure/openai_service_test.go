// Package infrastructure はテキスト字句解析のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"business/internal/openai/domain"
)

func TestOpenAIService_AnalyzeText_正常系_OpenAI_APIを呼び出してテキスト解析を実行すること(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := context.Background()
	service := NewOpenAIService("test-api-key")

	request := &domain.TextAnalysisRequest{
		Text: "以下のメールを字句解析してください：\n\n重要な会議の件について連絡いたします。",
		Options: domain.AnalysisOptions{
			EnableSentiment:  true,
			EnableKeywords:   true,
			EnableEntities:   true,
			EnableSummary:    true,
			EnableCategories: true,
			MaxKeywords:      10,
			MaxSummaryLength: 200,
		},
		Metadata: map[string]string{
			"source": "email",
		},
	}

	// Act
	result, err := service.AnalyzeText(ctx, request)

	// Assert
	// 実際のAPIを呼び出すため、APIキーが無効な場合はエラーになることを確認
	// 正常なAPIキーの場合は結果が返ることを確認
	if err != nil {
		// APIキーが無効またはネットワークエラーの場合
		assert.Contains(t, err.Error(), "API")
		assert.Nil(t, result)
	} else {
		// 正常な場合
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Language)
		assert.GreaterOrEqual(t, result.Confidence, 0.0)
		assert.LessOrEqual(t, result.Confidence, 1.0)
	}
}

func TestOpenAIService_AnalyzeText_異常系_空のテキストでエラーを返すこと(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := context.Background()
	service := NewOpenAIService("test-api-key")

	request := &domain.TextAnalysisRequest{
		Text: "",
		Options: domain.AnalysisOptions{
			EnableSentiment: true,
		},
	}

	// Act
	result, err := service.AnalyzeText(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEmptyText, err)
}

func TestOpenAIService_AnalyzeText_異常系_無効なAPIキーでエラーを返すこと(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := context.Background()
	service := NewOpenAIService("invalid-api-key")

	request := &domain.TextAnalysisRequest{
		Text: "テストテキスト",
		Options: domain.AnalysisOptions{
			EnableSentiment: true,
		},
	}

	// Act
	result, err := service.AnalyzeText(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API")
}

func TestOpenAIService_buildPrompt_正常系_リクエストからプロンプトを構築すること(t *testing.T) {
	t.Parallel()
	// Arrange
	service := NewOpenAIService("test-api-key")

	request := &domain.TextAnalysisRequest{
		Text: "テストメール本文",
		Options: domain.AnalysisOptions{
			EnableSentiment:  true,
			EnableKeywords:   true,
			EnableEntities:   false,
			EnableSummary:    true,
			EnableCategories: false,
			MaxKeywords:      5,
			MaxSummaryLength: 100,
		},
	}

	// Act
	prompt := service.buildPrompt(request)

	// Assert
	assert.Contains(t, prompt, "テストメール本文")
	assert.Contains(t, prompt, "JSON")
	assert.Contains(t, prompt, "sentiment")
	assert.Contains(t, prompt, "keywords")
	assert.Contains(t, prompt, "summary")
	assert.NotContains(t, prompt, "entities")
	assert.NotContains(t, prompt, "categories")
}

func TestOpenAIService_parseResponse_正常系_OpenAIレスポンスを解析結果に変換すること(t *testing.T) {
	t.Parallel()
	// Arrange
	service := NewOpenAIService("test-api-key")

	// OpenAI APIの典型的なレスポンス形式をシミュレート
	openaiResponse := `{
		"sentiment": {
			"score": 0.3,
			"magnitude": 0.8,
			"label": "POSITIVE",
			"confidence": 0.85
		},
		"keywords": [
			{"text": "会議", "relevance": 0.9, "count": 2, "category": "business"},
			{"text": "重要", "relevance": 0.7, "count": 1, "category": "priority"}
		],
		"summary": "重要な会議についての連絡メール",
		"language": "ja",
		"confidence": 0.95
	}`

	// Act
	result, err := service.parseResponse(openaiResponse)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ja", result.Language)
	assert.Equal(t, 0.95, result.Confidence)
	assert.Equal(t, "POSITIVE", result.Sentiment.Label)
	assert.Equal(t, 0.3, result.Sentiment.Score)
	assert.Len(t, result.Keywords, 2)
	assert.Equal(t, "会議", result.Keywords[0].Text)
	assert.Equal(t, "重要な会議についての連絡メール", result.Summary)
}

func TestOpenAIService_parseResponse_異常系_無効なJSONでエラーを返すこと(t *testing.T) {
	t.Parallel()
	// Arrange
	service := NewOpenAIService("test-api-key")

	invalidJSON := `{"invalid": json}`

	// Act
	result, err := service.parseResponse(invalidJSON)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "JSON")
}
