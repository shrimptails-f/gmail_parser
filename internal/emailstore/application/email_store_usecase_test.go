package application

import (
	cd "business/internal/common/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailStoreRepository の定義
type MockEmailStoreRepository struct {
	mock.Mock
}

func (m *MockEmailStoreRepository) SaveEmail(result cd.Email) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockEmailStoreRepository) GetEmailByGmailIds(gmailIds []string) ([]string, error) {
	args := m.Called(gmailIds)
	return args.Get(0).([]string), args.Error(1)
}

// テスト: SaveEmailAnalysisResult 成功時
func TestSaveEmailAnalysisResult_Success(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	email := cd.Email{GmailID: "test@gmail.com"}
	mockRepo.On("SaveEmail", email).Return(nil)

	err := usecase.SaveEmailAnalysisResult(email)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// テスト: SaveEmailAnalysisResult エラー時
func TestSaveEmailAnalysisResult_Error(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	email := cd.Email{GmailID: "test@gmail.com"}
	mockRepo.On("SaveEmail", email).Return(errors.New("db error"))

	err := usecase.SaveEmailAnalysisResult(email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "メール保存エラー")
	mockRepo.AssertExpectations(t)
}

// テスト: GetEmailByGmailIds 成功時
func TestGetEmailByGmailIds_Success(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	emailIds := []string{"gmail-id-1", "gmail-id-2"}
	expectedResult := []string{"gmail-id-1"}
	mockRepo.On("GetEmailByGmailIds", emailIds).Return(expectedResult, nil)

	result, err := usecase.GetEmailByGmailIds(emailIds)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

// テスト: GetEmailByGmailIds エラー時
func TestGetEmailByGmailIds_Error(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	emailIds := []string{"gmail-id-1"}
	mockRepo.On("GetEmailByGmailIds", emailIds).Return([]string{}, errors.New("db error"))

	result, err := usecase.GetEmailByGmailIds(emailIds)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "メール存在チェックエラー")
	mockRepo.AssertExpectations(t)
}

// テスト: GetEmailByGmailIds 空のリスト
func TestGetEmailByGmailIds_EmptyList(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	result, err := usecase.GetEmailByGmailIds([]string{})
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "メールIDが空です")
}
