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

func (m *MockEmailStoreRepository) CheckGmailIdExists(ctx context.Context, gmailId string) (*domain.Email, error) {
	args := m.Called(ctx, gmailId)
	return args.Get(0).(*domain.Email), args.Error(1)
}

func (m *MockEmailStoreRepository) EmailExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailStoreRepository) KeywordExists(word string) (bool, error) {
	args := m.Called(word)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailStoreRepository) PositionExists(ctx context.Context, word string) (bool, error) {
	args := m.Called(ctx, word)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailStoreRepository) WorkTypeExists(ctx context.Context, word string) (bool, error) {
	args := m.Called(ctx, word)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailStoreRepository) SaveEmailMultiple(ctx context.Context, result *openaidomain.EmailAnalysisMultipleResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockEmailStoreRepository) GetEmailByGmailId(ctx context.Context, gmail_id string) (*domain.Email, error) {
	args := m.Called(ctx, gmail_id)
	return args.Get(0).(*domain.Email), args.Error(1)
}

func TestEmailStoreUseCaseImpl_SaveEmailAnalysisResult(t *testing.T) {
	t.Parallel()
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
				GmailID:      "test-email-id",
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
				GmailID: "test-email-id",
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

func TestEmailStoreUseCaseImpl_CheckEmailExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         string
		expected      bool
		expectedError string
	}{
		{
			name: "正常系_メール存在",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("EmailExists", mock.Anything, "existing-email-id").Return(true, nil).Once()
			},
			input:         "existing-email-id",
			expected:      true,
			expectedError: "",
		},
		{
			name: "正常系_メール存在しない",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("EmailExists", mock.Anything, "non-existing-email-id").Return(false, nil).Once()
			},
			input:         "non-existing-email-id",
			expected:      false,
			expectedError: "",
		},
		{
			name: "異常系_空のメールID",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				// モックの設定なし（呼び出されない）
			},
			input:         "",
			expected:      false,
			expectedError: "メールIDが空です",
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
			exists, err := useCase.CheckGmailIdExists(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, exists)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailStoreUseCaseImpl_CheckKeywordExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         string
		expected      bool
		expectedError string
	}{
		{
			name: "正常系_キーワード存在",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("KeywordExists", "Go").Return(true, nil).Once()
			},
			input:         "Go",
			expected:      true,
			expectedError: "",
		},
		{
			name: "正常系_キーワード存在しない",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("KeywordExists", "Java").Return(false, nil).Once()
			},
			input:         "Java",
			expected:      false,
			expectedError: "",
		},
		{
			name: "異常系_空のキーワード",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				// モックの設定なし（呼び出されない）
			},
			input:         "",
			expected:      false,
			expectedError: "キーワードが空です",
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
			exists, err := useCase.CheckKeywordExists(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, exists)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailStoreUseCaseImpl_CheckPositionExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         string
		expected      bool
		expectedError string
	}{
		{
			name: "正常系_ポジション存在",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("PositionExists", mock.Anything, "PM").Return(true, nil).Once()
			},
			input:         "PM",
			expected:      true,
			expectedError: "",
		},
		{
			name: "正常系_ポジション存在しない",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("PositionExists", mock.Anything, "SE").Return(false, nil).Once()
			},
			input:         "SE",
			expected:      false,
			expectedError: "",
		},
		{
			name: "異常系_空のポジション",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				// モックの設定なし（呼び出されない）
			},
			input:         "",
			expected:      false,
			expectedError: "ポジションが空です",
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
			exists, err := useCase.CheckPositionExists(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, exists)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailStoreUseCaseImpl_CheckWorkTypeExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         string
		expected      bool
		expectedError string
	}{
		{
			name: "正常系_業務種別存在",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("WorkTypeExists", mock.Anything, "バックエンド開発").Return(true, nil).Once()
			},
			input:         "バックエンド開発",
			expected:      true,
			expectedError: "",
		},
		{
			name: "正常系_業務種別存在しない",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				mockRepo.On("WorkTypeExists", mock.Anything, "フロントエンド開発").Return(false, nil).Once()
			},
			input:         "フロントエンド開発",
			expected:      false,
			expectedError: "",
		},
		{
			name: "異常系_空の業務種別",
			setupMock: func(mockRepo *MockEmailStoreRepository) {
				// モックの設定なし（呼び出されない）
			},
			input:         "",
			expected:      false,
			expectedError: "業務種別が空です",
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
			exists, err := useCase.CheckWorkTypeExists(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, exists)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// intPtr はintのポインタを返すヘルパー関数です
func intPtr(i int) *int {
	return &i
}
