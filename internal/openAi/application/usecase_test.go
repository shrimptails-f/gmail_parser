package application

import (
	cd "business/internal/common/domain"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
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

	expected := []cd.Email{
		{

			GmailID:             "test-email-id-1",
			Summary:             "テスト案件名",
			Subject:             "テスト件名",
			ProjectName:         "テスト案件名",
			From:                "sender@example.com",
			FromEmail:           "sender@example.com",
			ReceivedDate:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Body:                "テスト本文",
			Category:            "案件",
			StartPeriod:         []string{"2024年4月", "2024年5月"},
			EndPeriod:           "2024年12月",
			WorkLocation:        "東京都",
			PriceFrom:           lo.ToPtr(500000),
			PriceTo:             lo.ToPtr(600000),
			Languages:           []string{"Go", "Python"},
			Frameworks:          []string{"Gin", "Django"},
			Positions:           []string{"PM", "SE"},
			WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
			RequiredSkillsMust:  []string{"Git", "Docker"},
			RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
			RemoteWorkCategory:  lo.ToPtr("フルリモート"),
			RemoteWorkFrequency: lo.ToPtr("週5日"),
		},
	}

	analyzeEmailBodyexpected := []cd.AnalysisResult{
		{
			MailCategory:        "案件",
			ProjectTitle:        "テスト案件名",
			StartPeriod:         []string{"2024年4月", "2024年5月"},
			EndPeriod:           "2024年12月",
			WorkLocation:        "東京都",
			PriceFrom:           lo.ToPtr(500000),
			PriceTo:             lo.ToPtr(600000),
			Languages:           []string{"Go", "Python"},
			Frameworks:          []string{"Gin", "Django"},
			Positions:           []string{"PM", "SE"},
			WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
			RequiredSkillsMust:  []string{"Git", "Docker"},
			RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
			RemoteWorkCategory:  lo.ToPtr("フルリモート"),
			RemoteWorkFrequency: lo.ToPtr("週5日"),
		},
	}

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
	mockAnalyzer.On("AnalyzeEmailBody", ctx, "PROMPT\n\nテスト本文").Return(analyzeEmailBodyexpected, nil)
	usecase := NewUseCase(mockAnalyzer, mockOS)

	input := []cd.BasicMessage{
		{
			ID:      expected[0].GmailID,
			Subject: expected[0].Subject,
			From:    expected[0].From,
			To:      []string{""},
			Date:    expected[0].ReceivedDate,
			Body:    expected[0].Body,
		},
	}
	actual, err := usecase.AnalyzeEmailContent(ctx, input)

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

	input := []cd.BasicMessage{}
	results, err := usecase.AnalyzeEmailContent(ctx, input)

	assert.Nil(t, results)
	assert.EqualError(t, err, "read error")
}
