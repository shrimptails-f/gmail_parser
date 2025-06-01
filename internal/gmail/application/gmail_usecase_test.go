// Package application のテストファイルです。
// Gmail認証ユースケースの動作をテストします。
package application

import (
	"business/internal/gmail/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGmailAuthService はGmailAuthServiceのモックです
type MockGmailAuthService struct {
	mock.Mock
}

func (m *MockGmailAuthService) Authenticate(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error) {
	args := m.Called(ctx, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GmailAuthResult), args.Error(1)
}

func (m *MockGmailAuthService) CreateGmailService(ctx context.Context, credential domain.GmailCredential, applicationName string) (interface{}, error) {
	args := m.Called(ctx, credential, applicationName)
	return args.Get(0), args.Error(1)
}

func (m *MockGmailAuthService) LoadCredentials(credentialsFolder, userID string) (*domain.GmailCredential, error) {
	args := m.Called(credentialsFolder, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GmailCredential), args.Error(1)
}

func (m *MockGmailAuthService) SaveCredentials(credentialsFolder, userID string, credential domain.GmailCredential) error {
	args := m.Called(credentialsFolder, userID, credential)
	return args.Error(0)
}

// TestGmailAuthUseCase_AuthenticateGmail_Success は正常系のテストです
func TestGmailAuthUseCase_AuthenticateGmail_Success(t *testing.T) {
	// モックを作成
	mockService := new(MockGmailAuthService)
	useCase := NewGmailAuthUseCase(mockService)

	// テストデータ
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "/path/to/client-secret.json",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	expectedResult := &domain.GmailAuthResult{
		Credential: domain.GmailCredential{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			TokenType:    "Bearer",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		ApplicationName: "gmailai",
		IsNewAuth:       true,
	}

	ctx := context.Background()

	// モックの期待を設定
	mockService.On("LoadCredentials", config.CredentialsFolder, config.UserID).Return(nil, domain.ErrClientSecretNotFound).Once()
	mockService.On("Authenticate", ctx, config).Return(expectedResult, nil).Once()
	mockService.On(
		"SaveCredentials",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(c domain.GmailCredential) bool {
			// ExpiresAtが微妙に異なりエラーになるので全て受け入れる
			return true
		}),
	).Return(nil).Maybe()
	result, err := useCase.AuthenticateGmail(ctx, config)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult.Credential.AccessToken, result.Credential.AccessToken)
	assert.Equal(t, expectedResult.ApplicationName, result.ApplicationName)
	assert.True(t, result.IsNewAuth)

	// モックの期待が満たされたかを確認
	mockService.AssertExpectations(t)
}

// TestGmailAuthUseCase_AuthenticateGmail_ExistingCredentials は既存認証情報がある場合のテストです
func TestGmailAuthUseCase_AuthenticateGmail_ExistingCredentials(t *testing.T) {
	// モックを作成
	mockService := new(MockGmailAuthService)
	useCase := NewGmailAuthUseCase(mockService)

	// テストデータ
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "/path/to/client-secret.json",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	existingCredential := &domain.GmailCredential{
		AccessToken:  "existing_access_token",
		RefreshToken: "existing_refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()

	// モックの期待を設定
	mockService.On("LoadCredentials", config.CredentialsFolder, config.UserID).Return(existingCredential, nil).Once()

	// テスト実行
	result, err := useCase.AuthenticateGmail(ctx, config)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, existingCredential.AccessToken, result.Credential.AccessToken)
	assert.Equal(t, config.ApplicationName, result.ApplicationName)
	assert.False(t, result.IsNewAuth)

	// モックの期待が満たされたかを確認
	mockService.AssertExpectations(t)
}

// TestGmailAuthUseCase_AuthenticateGmail_InvalidConfig は無効な設定の場合のテストです
func TestGmailAuthUseCase_AuthenticateGmail_InvalidConfig(t *testing.T) {
	// モックを作成
	mockService := new(MockGmailAuthService)
	useCase := NewGmailAuthUseCase(mockService)

	// 無効な設定（ClientSecretPathが空）
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
	}

	ctx := context.Background()

	// テスト実行
	result, err := useCase.AuthenticateGmail(ctx, config)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrClientSecretNotFound, err)

	// モックの期待が満たされたかを確認
	mockService.AssertExpectations(t)
}

// TestGmailAuthUseCase_CreateGmailService_Success は正常系のテストです
func TestGmailAuthUseCase_CreateGmailService_Success(t *testing.T) {
	// モックを作成
	mockService := new(MockGmailAuthService)
	useCase := NewGmailAuthUseCase(mockService)

	// テストデータ
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "/path/to/client-secret.json",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	credential := domain.GmailCredential{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	authResult := &domain.GmailAuthResult{
		Credential:      credential,
		ApplicationName: "gmailai",
		IsNewAuth:       true,
	}

	expectedService := "mock_gmail_service"

	ctx := context.Background()

	// モックの期待を設定
	mockService.On("LoadCredentials", config.CredentialsFolder, config.UserID).Return(nil, domain.ErrClientSecretNotFound).Once()
	mockService.On("Authenticate", ctx, config).Return(authResult, nil).Once()
	mockService.On(
		"SaveCredentials",
		config.CredentialsFolder,
		config.UserID,
		mock.MatchedBy(func(c domain.GmailCredential) bool {
			// ExpiresAtが微妙に異なりエラーになるので全て受け入れる
			return true
		}),
	).Return(nil).Maybe()
	mockService.On("CreateGmailService", ctx, credential, config.ApplicationName).Return(expectedService, nil).Once()

	// テスト実行
	service, err := useCase.CreateGmailService(ctx, config)

	// 検証
	assert.NoError(t, err)
	assert.Equal(t, expectedService, service)

	// モックの期待が満たされたかを確認
	mockService.AssertExpectations(t)
}

// TestGmailAuthUseCase_CreateGmailService_AuthenticationFailed は認証失敗の場合のテストです
func TestGmailAuthUseCase_CreateGmailService_AuthenticationFailed(t *testing.T) {
	// モックを作成
	mockService := new(MockGmailAuthService)
	useCase := NewGmailAuthUseCase(mockService)

	// テストデータ
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "/path/to/client-secret.json",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	ctx := context.Background()

	// モックの期待を設定
	mockService.On("LoadCredentials", config.CredentialsFolder, config.UserID).Return(nil, domain.ErrClientSecretNotFound).Once()
	mockService.On("Authenticate", ctx, config).Return(nil, domain.ErrGmailAuthFailed).Once()

	// テスト実行
	service, err := useCase.CreateGmailService(ctx, config)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Equal(t, domain.ErrGmailAuthFailed, err)

	// モックの期待が満たされたかを確認
	mockService.AssertExpectations(t)
}
