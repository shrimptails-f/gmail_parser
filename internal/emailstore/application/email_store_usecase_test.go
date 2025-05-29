// Package application はメール保存機能のアプリケーション層のテストを提供します。
package application

import (
	"business/internal/emailstore/domain"
	openaidomain "business/internal/openai/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailStoreRepository はEmailStoreRepositoryのモックです
type MockEmailStoreRepository struct {
	mock.Mock
}

func (m *MockEmailStoreRepository) SaveEmail(ctx context.Context, result *openaidomain.EmailAnalysisResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockEmailStoreRepository) GetEmailByID(ctx context.Context, id string) (*domain.Email, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Email), args.Error(1)
}

func (m *MockEmailStoreRepository) EmailExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestEmailStoreUseCaseImpl_SaveEmailAnalysisResult(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         *openaidomain.EmailAnalysisResult
		expectedError string
	}{
		{
			name: "正常系_新規メール保存成功",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("SaveEmail", mock.Anything, mock.AnythingOfType("*domain.EmailAnalysisResult")).Return(nil).Once()
			},
			input: &openaidomain.EmailAnalysisResult{
				ID:           "test-email-id",
				Subject:      "テスト件名",
				From:         "sender@example.com",
				FromEmail:    "sender@example.com",
				Date:         time.Now(),
				Body:         "テスト本文",
				MailCategory: "案件",
				StartPeriod:  []string{"2024年4月"},
				EndPeriod:    "2024年12月",
				WorkLocation: "東京都",
				PriceFrom:    intPtr(500000),
				PriceTo:      intPtr(600000),
				Languages:    []string{"Go", "Python"},
				Frameworks:   []string{"Gin", "Django"},
			},
			expectedError: "",
		},
		{
			name: "異常系_リポジトリエラー",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("SaveEmail", mock.Anything, mock.AnythingOfType("*domain.EmailAnalysisResult")).Return(domain.ErrInvalidEmailData).Once()
			},
			input: &openaidomain.EmailAnalysisResult{
				ID:      "test-email-id",
				Subject: "テスト件名",
			},
			expectedError: "メール保存エラー: 無効なメールデータです",
		},
		{
			name: "異常系_nilの分析結果",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				// モックの設定なし（呼び出されない）
			},
			input:         nil,
			expectedError: "分析結果がnilです",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockEmailStoreRepository)
			tt.setupMock(mockRepo)
			useCase := NewEmailStoreUseCase(mockRepo)
			ctx := context.Background()

			// Act
			err := useCase.SaveEmailAnalysisResult(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// intPtr はintのポインタを返すヘルパー関数です
func intPtr(i int) *int {
	return &i
}
