// Package application はテキスト字句解析のアプリケーション層のテストを提供します。
package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"business/internal/openai/domain"
)

// MockTextAnalysisService はテキスト解析サービスのモックです
type MockTextAnalysisService struct {
	mock.Mock
}

func (m *MockTextAnalysisService) AnalyzeText(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TextAnalysisResult), args.Error(1)
}

// MockPromptService はプロンプトサービスのモックです
type MockPromptService struct {
	mock.Mock
}

func (m *MockPromptService) LoadPrompt(filename string) (string, error) {
	args := m.Called(filename)
	return args.String(0), args.Error(1)
}

func (m *MockPromptService) SavePrompt(filename, content string) error {
	args := m.Called(filename, content)
	return args.Error(0)
}

func (m *MockPromptService) ListPrompts() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func TestTextAnalysisUseCase_AnalyzeText_正常系_テキストをAI解析すること(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTextAnalysisService := &MockTextAnalysisService{}
	mockPromptService := &MockPromptService{}

	useCase := NewTextAnalysisUseCase(mockTextAnalysisService, mockPromptService)

	text := "重要な会議の件について連絡いたします。"

	expectedResult := &domain.TextAnalysisResult{
		Summary:    "重要な会議についての連絡",
		Language:   "ja",
		Confidence: 0.95,
	}

	mockTextAnalysisService.On("AnalyzeText", ctx, mock.MatchedBy(func(req *domain.TextAnalysisRequest) bool {
		return req.Text == text &&
			req.Options.EnableSentiment
	})).Return(expectedResult, nil)

	// Act
	result, err := useCase.AnalyzeText(ctx, text)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockTextAnalysisService.AssertExpectations(t)
}

func TestTextAnalysisUseCase_AnalyzeText_異常系_AI解析失敗時にエラーを返すこと(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTextAnalysisService := &MockTextAnalysisService{}
	mockPromptService := &MockPromptService{}

	useCase := NewTextAnalysisUseCase(mockTextAnalysisService, mockPromptService)

	text := "テストメール"
	expectedError := domain.ErrAnalysisAPIError

	mockTextAnalysisService.On("AnalyzeText", ctx, mock.Anything).Return(nil, expectedError)

	// Act
	result, err := useCase.AnalyzeText(ctx, text)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "テキスト解析エラー")
	assert.Contains(t, err.Error(), expectedError.Error())
	mockTextAnalysisService.AssertExpectations(t)
}

func TestTextAnalysisUseCase_AnalyzeText_異常系_空のテキストでエラーを返すこと(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTextAnalysisService := &MockTextAnalysisService{}
	mockPromptService := &MockPromptService{}

	useCase := NewTextAnalysisUseCase(mockTextAnalysisService, mockPromptService)

	// Act
	result, err := useCase.AnalyzeText(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrEmptyText, err)
}

func TestTextAnalysisUseCase_AnalyzeTextWithOptions_正常系_オプション付きでテキストを解析すること(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTextAnalysisService := &MockTextAnalysisService{}
	mockPromptService := &MockPromptService{}

	useCase := NewTextAnalysisUseCase(mockTextAnalysisService, mockPromptService)

	request := &domain.TextAnalysisRequest{
		Text: "テストメール",
		Options: domain.AnalysisOptions{
			EnableSentiment: true,
			EnableKeywords:  true,
		},
	}

	expectedResult := &domain.TextAnalysisResult{
		Summary:    "テストメールの要約",
		Language:   "ja",
		Confidence: 0.95,
	}

	mockTextAnalysisService.On("AnalyzeText", ctx, request).Return(expectedResult, nil)

	// Act
	result, err := useCase.AnalyzeTextWithOptions(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockTextAnalysisService.AssertExpectations(t)
}

func TestTextAnalysisUseCase_DisplayAnalysisResult_正常系_解析結果をターミナルに表示すること(t *testing.T) {
	// Arrange
	mockTextAnalysisService := &MockTextAnalysisService{}
	mockPromptService := &MockPromptService{}

	useCase := NewTextAnalysisUseCase(mockTextAnalysisService, mockPromptService)

	result := &domain.TextAnalysisResult{
		MessageID:  "test-message-id",
		Subject:    "会議の件",
		Summary:    "重要な会議についての連絡",
		Language:   "ja",
		Confidence: 0.95,
		Sentiment: domain.SentimentAnalysis{
			Score:      0.3,
			Label:      "POSITIVE",
			Confidence: 0.8,
		},
		Keywords: []domain.Keyword{
			{Text: "会議", Relevance: 0.9, Count: 2},
			{Text: "重要", Relevance: 0.7, Count: 1},
		},
	}

	// Act & Assert
	// この関数は標準出力に表示するため、エラーが発生しないことを確認
	err := useCase.DisplayAnalysisResult(result)
	assert.NoError(t, err)
}
