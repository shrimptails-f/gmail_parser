package application

import (
	cd "business/internal/common/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モック: oswrapper
type mockOsWrapper struct {
	ReadFileFunc func(path string) (string, error)
	GetEnvFunc   func(key string) string
}

func (m *mockOsWrapper) ReadFile(path string) (string, error) {
	return m.ReadFileFunc(path)
}

func (m *mockOsWrapper) GetEnv(key string) string {
	return m.GetEnvFunc(key)
}

type mockAnalyzer struct {
	mock.Mock
}

func (m *mockAnalyzer) AnalyzeEmailBody(ctx context.Context, prompt string) ([]cd.AnalysisResult, error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).([]cd.AnalysisResult), args.Error(1)
}

// テスト関数

func TestAnalyzeEmailContent_Success(t *testing.T) {
	ctx := context.Background()

	mockOS := &mockOsWrapper{
		ReadFileFunc: func(path string) (string, error) {
			assert.Equal(t, "/data/prompts/text_analysis_prompt.txt", path)
			return "PROMPT", nil
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	mockAnalyzer := new(mockAnalyzer)

	expected := []cd.AnalysisResult{
		{MailCategory: "案件", ProjectTitle: "Go案件"},
	}
	mockAnalyzer.On("AnalyzeEmailBody", ctx, "PROMPT\n\nMAILBODY").Return(expected, nil)

	usecase := NewUseCase(mockAnalyzer, mockOS)

	actual, err := usecase.AnalyzeEmailContent(ctx, "MAILBODY")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	mockAnalyzer.AssertExpectations(t)
}

func TestAnalyzeEmailContent_ReadFileError(t *testing.T) {
	ctx := context.Background()

	mockOS := &mockOsWrapper{
		ReadFileFunc: func(path string) (string, error) {
			return "", errors.New("read error")
		},
		GetEnvFunc: func(key string) string {
			return ""
		},
	}

	mockAnalyzer := new(mockAnalyzer)

	usecase := NewUseCase(mockAnalyzer, mockOS)

	results, err := usecase.AnalyzeEmailContent(ctx, "anything")

	assert.Nil(t, results)
	assert.EqualError(t, err, "read error")
}
