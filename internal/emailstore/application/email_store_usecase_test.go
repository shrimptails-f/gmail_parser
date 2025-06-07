package application

import (
	cd "business/internal/common/domain"
	r "business/internal/emailstore/infrastructure"
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

func (m *MockEmailStoreRepository) GetEmailByGmailId(gmailId string) (r.Email, error) {
	args := m.Called(gmailId)
	return args.Get(0).(r.Email), args.Error(1)
}

func (m *MockEmailStoreRepository) EmailExists(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
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

// テスト: CheckGmailIdExists 成功 (存在)
func TestCheckGmailIdExists_Exists(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	expectedEmail := r.Email{ID: 123}
	mockRepo.On("GetEmailByGmailId", "gmail-id-123").Return(expectedEmail, nil)

	exists, err := usecase.CheckGmailIdExists("gmail-id-123")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

// テスト: CheckGmailIdExists 成功 (存在しない)
func TestCheckGmailIdExists_NotExists(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	expectedEmail := r.Email{}
	mockRepo.On("GetEmailByGmailId", "non-existent-id").Return(expectedEmail, nil)

	exists, err := usecase.CheckGmailIdExists("non-existent-id")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

// テスト: CheckGmailIdExists エラー発生
func TestCheckGmailIdExists_Error(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	mockRepo.On("GetEmailByGmailId", "error-id").Return(r.Email{}, errors.New("db error"))

	exists, err := usecase.CheckGmailIdExists("error-id")
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Contains(t, err.Error(), "メール存在チェックエラー")
	mockRepo.AssertExpectations(t)
}

// テスト: CheckGmailIdExists 空文字
func TestCheckGmailIdExists_Empty(t *testing.T) {
	mockRepo := new(MockEmailStoreRepository)
	usecase := NewEmailStoreUseCase(mockRepo)

	exists, err := usecase.CheckGmailIdExists("")
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Contains(t, err.Error(), "メールIDが空です")
}
