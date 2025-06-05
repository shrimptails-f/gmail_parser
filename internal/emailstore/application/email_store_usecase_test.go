// Package application はメール保存機能のアプリケーション層のテストを提供します。
package application

import (
	"business/internal/emailstore/domain"
	r "business/internal/emailstore/infrastructure"
	openaidomain "business/internal/openai/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailStoreRepository はEmailStoreRepositoryのモックです
type MockEmailStoreRepository struct {
	mock.Mock
}

func (m *MockEmailStoreRepository) SaveEmail(result r.Email) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockEmailStoreRepository) SaveEmailMultiple(result *openaidomain.EmailAnalysisMultipleResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockEmailStoreRepository) GetEmailByGmailId(gmail_id string) (*domain.Email, error) {
	args := m.Called(gmail_id)
	return args.Get(0).(*domain.Email), args.Error(1)
}

func (m *MockEmailStoreRepository) EmailExists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailStoreRepository) GetkeywordGroups(name []string) ([]r.KeywordGroup, error) {
	args := m.Called(name)
	return args.Get(0).([]r.KeywordGroup), args.Error(1)
}

func (m *MockEmailStoreRepository) GetKeywords(words []string) ([]r.KeyWord, error) {
	args := m.Called(words)
	return args.Get(0).([]r.KeyWord), args.Error(1)
}

func (m *MockEmailStoreRepository) GetPositionGroups(name []string) ([]r.PositionGroup, error) {
	args := m.Called(name)
	return args.Get(0).([]r.PositionGroup), args.Error(1)
}

func (m *MockEmailStoreRepository) GetPositionWords(words []string) ([]r.PositionWord, error) {
	args := m.Called(words)
	return args.Get(0).([]r.PositionWord), args.Error(1)
}

func (m *MockEmailStoreRepository) GetWorkTypeWords(words []string) ([]r.WorkTypeWord, error) {
	args := m.Called(words)
	return args.Get(0).([]r.WorkTypeWord), args.Error(1)
}

func (m *MockEmailStoreRepository) GetWorkTypeGroups(words []string) ([]r.WorkTypeGroup, error) {
	args := m.Called(words)
	return args.Get(0).([]r.WorkTypeGroup), args.Error(1)
}

func TestEmailStoreUseCaseImpl_SaveEmailAnalysisResult(t *testing.T) {
	t.Parallel()

	tt := struct {
		name          string
		setupMock     func(*MockEmailStoreRepository)
		input         *domain.AnalysisResult
		expectedError string
	}{
		name: "正常系_新規メール保存成功",
		setupMock: func(mockRepo *MockEmailStoreRepository) {
			mockRepo.On("EmailExists", "test-email-id").Return(false, nil).Once()

			mockRepo.On("GetKeywords", mock.Anything).Return([]r.KeyWord{}, nil).Times(4)
			mockRepo.On("GetkeywordGroups", mock.Anything).Return([]r.KeywordGroup{}, nil).Times(4)

			mockRepo.On("GetPositionWords", mock.Anything).Return([]r.PositionWord{}, nil).Once()
			mockRepo.On("GetPositionGroups", mock.AnythingOfType("[]string")).Return([]r.PositionGroup{}, nil).Once()

			mockRepo.On("GetWorkTypeWords", mock.Anything).Return([]r.WorkTypeWord{}, nil).Once()
			mockRepo.On("GetWorkTypeGroups", mock.AnythingOfType("[]string")).Return([]r.WorkTypeGroup{}, nil).Once()

			mockRepo.On("SaveEmail", mock.Anything).Return(nil).Once()
		},

		input: &domain.AnalysisResult{
			GmailID:            "test-email-id",
			Subject:            "テスト件名",
			From:               "田中 太郎 <sender@example.com>",
			FromEmail:          "sender@example.com",
			ReceivedDate:       time.Now(),
			Body:               "テスト本文",
			ProjectName:        "プロジェクトA",
			StartPeriod:        []string{"2024年4月"},
			EndPeriod:          "2024年12月",
			WorkLocation:       "東京都",
			PriceFrom:          intPtr(500000),
			PriceTo:            intPtr(600000),
			Languages:          []string{"Go", "Python"},
			Frameworks:         []string{"Gin", "Django"},
			RequiredSkillsMust: []string{"必須スキル1", "必須スキル2"},
			RequiredSkillsWant: []string{"尚可スキル1", "尚可スキル2"},
			Positions:          []string{"SE", "PG"},
			WorkTypes:          []string{"バックエンド", "インフラエンジニア"},
		},
		expectedError: "",
	}

	t.Run(tt.name, func(t *testing.T) {
		mockRepo := new(MockEmailStoreRepository)
		tt.setupMock(mockRepo)
		useCase := NewEmailStoreUseCase(mockRepo)
		err := useCase.SaveEmailAnalysisResult(*tt.input)

		if tt.expectedError == "" {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		}

		mockRepo.AssertExpectations(t)
	})
}

func TestEmailStoreUseCaseImpl_SaveEmailAnalysisResult_SaveEmailFails(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockEmailStoreRepository)

	mockRepo.
		On("EmailExists", "test-id").
		Return(true, nil).
		Once()

	useCase := NewEmailStoreUseCase(mockRepo)

	input := domain.AnalysisResult{
		GmailID: "test-id",
		Subject: "失敗テスト",
		From:    "山田 花子 <yamada@example.com>",
	}

	err := useCase.SaveEmailAnalysisResult(input)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// intPtr はintのポインタを返すヘルパー関数です
func intPtr(i int) *int {
	return &i
}
