// Package application はメール分析のアプリケーション層のテストを提供します。
package application

import (
	gmail "business/internal/gmail/domain"
	"business/internal/openai/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailAnalysisService はメール分析サービスのモックです
type MockEmailAnalysisService struct {
	mock.Mock
}

func (m *MockEmailAnalysisService) AnalyzeEmail(ctx context.Context, request *domain.EmailAnalysisRequest) (*domain.EmailAnalysisResult, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*domain.EmailAnalysisResult), args.Error(1)
}

func TestEmailAnalysisUseCaseImpl_AnalyzeEmailContent(t *testing.T) {
	tests := []struct {
		name          string
		message       *gmail.GmailMessage
		mockResult    *domain.EmailAnalysisResult
		mockError     error
		expectedError string
	}{
		{
			name: "正常系_メール分析成功",
			message: &gmail.GmailMessage{
				ID:      "test-id",
				Subject: "テスト件名",
				From:    "Test User <test@example.com>",
				Date:    time.Now(),
				Body:    "テストメール本文",
			},
			mockResult: &domain.EmailAnalysisResult{
				PriceFrom:           intPtr(700000),
				PriceTo:             intPtr(900000),
				MailCategory:        "案件",
				StartPeriod:         []string{"2025/06/01"},
				EndPeriod:           "~長期",
				WorkLocation:        "東京都",
				Languages:           []string{"Go", "TypeScript"},
				Frameworks:          []string{"React", "Vue.js"},
				Positions:           []string{"バックエンドエンジニア", "フロントエンドエンジニア"},
				WorkTypes:           []string{"設計", "実装", "テスト"},
				RequiredSkillsMust:  []string{"Go経験3年以上", "TypeScript経験2年以上"},
				RequiredSkillsWant:  []string{"AWS経験", "Docker経験"},
				RemoteWorkCategory:  "フルリモート可",
				RemoteWorkFrequency: stringPtr("週5日"),
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name: "異常系_空のメール本文",
			message: &gmail.GmailMessage{
				ID:      "test-id",
				Subject: "テスト件名",
				From:    "Test User <test@example.com>",
				Date:    time.Now(),
				Body:    "",
			},
			mockResult:    nil,
			mockError:     nil,
			expectedError: "分析対象のメール本文が空です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスを作成
			mockService := new(MockEmailAnalysisService)

			// 空のメール本文でない場合のみモックの期待を設定
			if tt.message.Body != "" {
				mockService.On("AnalyzeEmail", mock.Anything, mock.AnythingOfType("*domain.EmailAnalysisRequest")).
					Return(tt.mockResult, tt.mockError)
			}

			// ユースケースを作成
			useCase := NewEmailAnalysisUseCase(mockService)

			// テスト実行
			result, err := useCase.AnalyzeEmailContent(context.Background(), tt.message)

			// 結果検証
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Gmail情報の検証
				assert.Equal(t, tt.message.ID, result.GmailID)
				assert.Equal(t, tt.message.Subject, result.Subject)
				assert.Equal(t, tt.message.From, result.From)
				assert.Equal(t, "test@example.com", result.FromEmail)
				assert.Equal(t, tt.message.Date, result.Date)
				assert.Equal(t, tt.message.Body, result.Body)
				// AI分析結果の検証
				assert.Equal(t, tt.mockResult.MailCategory, result.MailCategory)
				assert.Equal(t, tt.mockResult.StartPeriod, result.StartPeriod)
				assert.Equal(t, tt.mockResult.EndPeriod, result.EndPeriod)
				assert.Equal(t, tt.mockResult.WorkLocation, result.WorkLocation)
				assert.Equal(t, tt.mockResult.PriceFrom, result.PriceFrom)
				assert.Equal(t, tt.mockResult.PriceTo, result.PriceTo)
				assert.Equal(t, tt.mockResult.Languages, result.Languages)
				assert.Equal(t, tt.mockResult.Frameworks, result.Frameworks)
				assert.Equal(t, tt.mockResult.Positions, result.Positions)
				assert.Equal(t, tt.mockResult.WorkTypes, result.WorkTypes)
				assert.Equal(t, tt.mockResult.RequiredSkillsMust, result.RequiredSkillsMust)
				assert.Equal(t, tt.mockResult.RequiredSkillsWant, result.RequiredSkillsWant)
				assert.Equal(t, tt.mockResult.RemoteWorkCategory, result.RemoteWorkCategory)
				assert.Equal(t, tt.mockResult.RemoteWorkFrequency, result.RemoteWorkFrequency)
			}

			// モックの期待が満たされたか確認
			mockService.AssertExpectations(t)
		})
	}
}

func TestEmailAnalysisUseCaseImpl_DisplayEmailAnalysisResult(t *testing.T) {
	tests := []struct {
		name          string
		result        *domain.EmailAnalysisResult
		expectedError string
	}{
		{
			name: "正常系_結果表示成功",
			result: &domain.EmailAnalysisResult{
				Subject:             "テスト件名",
				From:                "Test User",
				FromEmail:           "test@example.com",
				Date:                time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
				Body:                "テストメール本文",
				MailCategory:        "案件",
				StartPeriod:         []string{"2025/06/01"},
				EndPeriod:           "~長期",
				WorkLocation:        "東京都",
				PriceFrom:           intPtr(800000),
				PriceTo:             intPtr(900000),
				Languages:           []string{"TypeScript", "JavaScript"},
				Frameworks:          []string{"React"},
				Positions:           []string{"フロントエンドエンジニア"},
				WorkTypes:           []string{"実装", "テスト"},
				RequiredSkillsMust:  []string{"React経験3年以上"},
				RequiredSkillsWant:  []string{"TypeScript経験"},
				RemoteWorkCategory:  "一部リモート可",
				RemoteWorkFrequency: stringPtr("週2-3日"),
			},
			expectedError: "",
		},
		{
			name:          "異常系_結果がnil",
			result:        nil,
			expectedError: "分析結果がnilです",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスを作成（DisplayEmailAnalysisResultでは使用されない）
			mockService := new(MockEmailAnalysisService)

			// ユースケースを作成
			useCase := NewEmailAnalysisUseCase(mockService)

			// テスト実行
			err := useCase.DisplayEmailAnalysisResult(tt.result)

			// 結果検証
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_extractEmailAddress(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		expected string
	}{
		{
			name:     "正常系_名前付きメールアドレス",
			from:     "Test User <test@example.com>",
			expected: "test@example.com",
		},
		{
			name:     "正常系_メールアドレスのみ",
			from:     "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "異常系_メールアドレスなし",
			from:     "Test User",
			expected: "",
		},
		{
			name:     "異常系_不正な形式",
			from:     "Test User <invalid>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractEmailAddress(tt.from)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_truncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "正常系_切り詰めなし",
			text:     "短いテキスト",
			maxLen:   20,
			expected: "短いテキスト",
		},
		{
			name:     "正常系_切り詰めあり",
			text:     "これは非常に長いテキストです",
			maxLen:   10,
			expected: "これは非常に長いテキ...",
		},
		{
			name:     "正常系_空文字",
			text:     "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_formatStringArray(t *testing.T) {
	tests := []struct {
		name     string
		arr      []string
		expected string
	}{
		{
			name:     "正常系_複数要素",
			arr:      []string{"TypeScript", "JavaScript", "PHP"},
			expected: `["TypeScript", "JavaScript", "PHP"]`,
		},
		{
			name:     "正常系_単一要素",
			arr:      []string{"React"},
			expected: `["React"]`,
		},
		{
			name:     "正常系_空配列",
			arr:      []string{},
			expected: "[]",
		},
		{
			name:     "正常系_nil配列",
			arr:      nil,
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStringArray(tt.arr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// intPtr はintのポインタを返すヘルパー関数です
func intPtr(i int) *int {
	return &i
}

// stringPtr はstringのポインタを返すヘルパー関数です
func stringPtr(s string) *string {
	return &s
}
