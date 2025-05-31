// Package application は認証機能のアプリケーション層を提供します。
// このファイルはGmailメッセージ取得ユースケースのテストを定義します。
package application

import (
	"business/internal/gmail/domain"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGmailMessageService はGmailMessageServiceのモックです
type MockGmailMessageService struct {
	mock.Mock
}

func (m *MockGmailMessageService) GetMessages(ctx context.Context, credential domain.GmailCredential, applicationName string, maxResults int64) ([]domain.GmailMessage, error) {
	args := m.Called(ctx, credential, applicationName, maxResults)
	return args.Get(0).([]domain.GmailMessage), args.Error(1)
}

func (m *MockGmailMessageService) GetMessage(ctx context.Context, credential domain.GmailCredential, applicationName string, messageID string) (*domain.GmailMessage, error) {
	args := m.Called(ctx, credential, applicationName, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GmailMessage), args.Error(1)
}

func (m *MockGmailMessageService) GetLabels(ctx context.Context, credential domain.GmailCredential, applicationName string) ([]domain.GmailInfo, error) {
	args := m.Called(ctx, credential, applicationName)
	return args.Get(0).([]domain.GmailInfo), args.Error(1)
}

func (m *MockGmailMessageService) GetMessagesByLabel(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, maxResults int64) ([]domain.GmailMessage, error) {
	args := m.Called(ctx, credential, applicationName, labelID, maxResults)
	return args.Get(0).([]domain.GmailMessage), args.Error(1)
}

func (m *MockGmailMessageService) GetMessagesByLabelWithPagination(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, maxResults int64, pageToken string) ([]domain.GmailMessage, string, error) {
	args := m.Called(ctx, credential, applicationName, labelID, maxResults, pageToken)
	return args.Get(0).([]domain.GmailMessage), args.String(1), args.Error(2)
}

func (m *MockGmailMessageService) GetMessagesByLabelAndDate(ctx context.Context, credential domain.GmailCredential, applicationName string, labelID string, afterDate time.Time, maxResults int64, pageToken string) ([]domain.GmailMessage, string, error) {
	args := m.Called(ctx, credential, applicationName, labelID, afterDate, maxResults, pageToken)
	return args.Get(0).([]domain.GmailMessage), args.String(1), args.Error(2)
}

func TestGmailMessageUseCaseImpl_GetMessages(t *testing.T) {
	tests := []struct {
		name           string
		config         domain.GmailAuthConfig
		maxResults     int64
		setupMock      func(*MockGmailMessageService)
		expectedResult []domain.GmailMessage
		expectedError  string
	}{
		{
			name: "正常系_メッセージ一覧取得成功",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			maxResults: 10,
			setupMock: func(mockService *MockGmailMessageService) {
				expectedMessages := []domain.GmailMessage{
					{
						ID:      "msg1",
						Subject: "テストメール1",
						From:    "sender1@example.com",
						To:      []string{"recipient@example.com"},
						Date:    time.Now(),
						Body:    "テストメール1の本文",
					},
					{
						ID:      "msg2",
						Subject: "テストメール2",
						From:    "sender2@example.com",
						To:      []string{"recipient@example.com"},
						Date:    time.Now(),
						Body:    "テストメール2の本文",
					},
				}
				mockService.On("GetMessages", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", int64(10)).Return(expectedMessages, nil)
			},
			expectedResult: []domain.GmailMessage{
				{
					ID:      "msg1",
					Subject: "テストメール1",
					From:    "sender1@example.com",
					To:      []string{"recipient@example.com"},
					Body:    "テストメール1の本文",
				},
				{
					ID:      "msg2",
					Subject: "テストメール2",
					From:    "sender2@example.com",
					To:      []string{"recipient@example.com"},
					Body:    "テストメール2の本文",
				},
			},
		},
		{
			name: "異常系_無効な設定",
			config: domain.GmailAuthConfig{
				ClientSecretPath: "", // 無効な設定
			},
			maxResults:    10,
			setupMock:     func(mockService *MockGmailMessageService) {},
			expectedError: "client-secret.jsonファイルが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockGmailAuthService := &MockGmailAuthService{}
			mockGmailMessageService := &MockGmailMessageService{}
			tt.setupMock(mockGmailMessageService)

			// 認証が必要な場合のモック設定
			if tt.expectedError == "" {
				credential := domain.GmailCredential{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "Bearer",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				mockGmailAuthService.On("LoadCredentials", tt.config.CredentialsFolder, "user").Return(&credential, nil)
			}

			useCase := NewGmailMessageUseCase(mockGmailAuthService, mockGmailMessageService)

			// Act
			result, err := useCase.GetMessages(context.Background(), tt.config, tt.maxResults)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
				for i, expectedMsg := range tt.expectedResult {
					assert.Equal(t, expectedMsg.ID, result[i].ID)
					assert.Equal(t, expectedMsg.Subject, result[i].Subject)
					assert.Equal(t, expectedMsg.From, result[i].From)
					assert.Equal(t, expectedMsg.To, result[i].To)
					assert.Equal(t, expectedMsg.Body, result[i].Body)
				}
			}

			// モックの期待が満たされたかチェック
			mockGmailAuthService.AssertExpectations(t)
			mockGmailMessageService.AssertExpectations(t)
		})
	}
}

func TestGmailMessageUseCaseImpl_GetAllMessagesByLabelPathFromToday(t *testing.T) {
	tests := []struct {
		name           string
		config         domain.GmailAuthConfig
		labelPath      string
		maxResults     int64
		setupMock      func(*MockGmailMessageService)
		expectedResult []domain.GmailMessage
		expectedError  string
	}{
		{
			name: "正常系_当日メッセージ全件取得成功_1ページ",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			labelPath:  "営業/案件",
			maxResults: 10,
			setupMock: func(mockService *MockGmailMessageService) {
				// ラベル一覧を返すモック
				labels := []domain.GmailInfo{
					{ID: "Label_1", Name: "営業", Type: "user"},
					{ID: "Label_2", Name: "営業/案件", Type: "user"},
					{ID: "Label_3", Name: "プライベート", Type: "user"},
				}
				mockService.On("GetLabels", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app").Return(labels, nil)

				// 当日のメッセージを返すモック（1ページ目で終了）
				expectedMessages := []domain.GmailMessage{
					{
						ID:      "msg1",
						Subject: "今日の案件: TypeScript開発",
						From:    "client@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "今日のTypeScript開発の案件です",
					},
					{
						ID:      "msg2",
						Subject: "今日の案件: Go開発",
						From:    "client2@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "今日のGo開発の案件です",
					},
				}
				mockService.On("GetMessagesByLabelAndDate", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", "Label_2", mock.AnythingOfType("time.Time"), int64(10), "").Return(expectedMessages, "", nil)
			},
			expectedResult: []domain.GmailMessage{
				{
					ID:      "msg1",
					Subject: "今日の案件: TypeScript開発",
					From:    "client@example.com",
					To:      []string{"sales@company.com"},
					Body:    "今日のTypeScript開発の案件です",
				},
				{
					ID:      "msg2",
					Subject: "今日の案件: Go開発",
					From:    "client2@example.com",
					To:      []string{"sales@company.com"},
					Body:    "今日のGo開発の案件です",
				},
			},
		},
		{
			name: "正常系_当日メッセージ全件取得成功_複数ページ",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			labelPath:  "営業/案件",
			maxResults: 2,
			setupMock: func(mockService *MockGmailMessageService) {
				// ラベル一覧を返すモック
				labels := []domain.GmailInfo{
					{ID: "Label_1", Name: "営業", Type: "user"},
					{ID: "Label_2", Name: "営業/案件", Type: "user"},
				}
				mockService.On("GetLabels", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app").Return(labels, nil)

				// 1ページ目のメッセージ
				firstPageMessages := []domain.GmailMessage{
					{
						ID:      "msg1",
						Subject: "案件1",
						From:    "client1@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "案件1の内容",
					},
					{
						ID:      "msg2",
						Subject: "案件2",
						From:    "client2@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "案件2の内容",
					},
				}
				mockService.On("GetMessagesByLabelAndDate", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", "Label_2", mock.AnythingOfType("time.Time"), int64(2), "").Return(firstPageMessages, "next_page_token", nil)

				// 2ページ目のメッセージ（最後のページ）
				secondPageMessages := []domain.GmailMessage{
					{
						ID:      "msg3",
						Subject: "案件3",
						From:    "client3@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "案件3の内容",
					},
				}
				mockService.On("GetMessagesByLabelAndDate", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", "Label_2", mock.AnythingOfType("time.Time"), int64(2), "next_page_token").Return(secondPageMessages, "", nil)
			},
			expectedResult: []domain.GmailMessage{
				{
					ID:      "msg1",
					Subject: "案件1",
					From:    "client1@example.com",
					To:      []string{"sales@company.com"},
					Body:    "案件1の内容",
				},
				{
					ID:      "msg2",
					Subject: "案件2",
					From:    "client2@example.com",
					To:      []string{"sales@company.com"},
					Body:    "案件2の内容",
				},
				{
					ID:      "msg3",
					Subject: "案件3",
					From:    "client3@example.com",
					To:      []string{"sales@company.com"},
					Body:    "案件3の内容",
				},
			},
		},
		{
			name: "異常系_ラベルが見つからない",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			labelPath:  "存在しないラベル",
			maxResults: 10,
			setupMock: func(mockService *MockGmailMessageService) {
				labels := []domain.GmailInfo{
					{ID: "Label_1", Name: "営業", Type: "user"},
					{ID: "Label_2", Name: "営業/案件", Type: "user"},
				}
				mockService.On("GetLabels", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app").Return(labels, nil)
			},
			expectedError: "指定されたラベル '存在しないラベル' が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockGmailAuthService := &MockGmailAuthService{}
			mockGmailMessageService := &MockGmailMessageService{}
			tt.setupMock(mockGmailMessageService)

			// 認証が必要な場合のモック設定
			if tt.expectedError == "" || !strings.Contains(tt.expectedError, "無効な設定") {
				credential := domain.GmailCredential{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "Bearer",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				mockGmailAuthService.On("LoadCredentials", tt.config.CredentialsFolder, "user").Return(&credential, nil)
			}

			useCase := NewGmailMessageUseCase(mockGmailAuthService, mockGmailMessageService)

			// Act
			result, err := useCase.GetAllMessagesByLabelPathFromToday(context.Background(), tt.config, tt.labelPath, tt.maxResults)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
				for i, expectedMsg := range tt.expectedResult {
					assert.Equal(t, expectedMsg.ID, result[i].ID)
					assert.Equal(t, expectedMsg.Subject, result[i].Subject)
					assert.Equal(t, expectedMsg.From, result[i].From)
					assert.Equal(t, expectedMsg.To, result[i].To)
					assert.Equal(t, expectedMsg.Body, result[i].Body)
				}
			}

			// モックの期待が満たされたかチェック
			mockGmailAuthService.AssertExpectations(t)
			mockGmailMessageService.AssertExpectations(t)
		})
	}
}

func TestGmailMessageUseCaseImpl_GetMessagesByLabelPath(t *testing.T) {
	tests := []struct {
		name           string
		config         domain.GmailAuthConfig
		labelPath      string
		maxResults     int64
		setupMock      func(*MockGmailMessageService)
		expectedResult []domain.GmailMessage
		expectedError  string
	}{
		{
			name: "正常系_ラベル指定メッセージ取得成功",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			labelPath:  "営業/案件",
			maxResults: 5,
			setupMock: func(mockService *MockGmailMessageService) {
				// ラベル一覧を返すモック
				labels := []domain.GmailInfo{
					{ID: "Label_1", Name: "営業", Type: "user"},
					{ID: "Label_2", Name: "営業/案件", Type: "user"},
					{ID: "Label_3", Name: "プライベート", Type: "user"},
				}
				mockService.On("GetLabels", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app").Return(labels, nil)

				// 指定されたラベルのメッセージを返すモック
				expectedMessages := []domain.GmailMessage{
					{
						ID:      "msg1",
						Subject: "案件: TypeScript開発",
						From:    "client@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "TypeScript開発の案件です",
					},
					{
						ID:      "msg2",
						Subject: "案件: Go開発",
						From:    "client2@example.com",
						To:      []string{"sales@company.com"},
						Date:    time.Now(),
						Body:    "Go開発の案件です",
					},
				}
				mockService.On("GetMessagesByLabel", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", "Label_2", int64(5)).Return(expectedMessages, nil)
			},
			expectedResult: []domain.GmailMessage{
				{
					ID:      "msg1",
					Subject: "案件: TypeScript開発",
					From:    "client@example.com",
					To:      []string{"sales@company.com"},
					Body:    "TypeScript開発の案件です",
				},
				{
					ID:      "msg2",
					Subject: "案件: Go開発",
					From:    "client2@example.com",
					To:      []string{"sales@company.com"},
					Body:    "Go開発の案件です",
				},
			},
		},
		{
			name: "異常系_ラベルが見つからない",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			labelPath:  "存在しないラベル",
			maxResults: 5,
			setupMock: func(mockService *MockGmailMessageService) {
				labels := []domain.GmailInfo{
					{ID: "Label_1", Name: "営業", Type: "user"},
					{ID: "Label_2", Name: "営業/案件", Type: "user"},
				}
				mockService.On("GetLabels", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app").Return(labels, nil)
			},
			expectedError: "指定されたラベル '存在しないラベル' が見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockGmailAuthService := &MockGmailAuthService{}
			mockGmailMessageService := &MockGmailMessageService{}
			tt.setupMock(mockGmailMessageService)

			// 認証が必要な場合のモック設定
			if tt.expectedError == "" || !strings.Contains(tt.expectedError, "無効な設定") {
				credential := domain.GmailCredential{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "Bearer",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				mockGmailAuthService.On("LoadCredentials", tt.config.CredentialsFolder, "user").Return(&credential, nil)
			}

			useCase := NewGmailMessageUseCase(mockGmailAuthService, mockGmailMessageService)

			// Act
			result, err := useCase.GetMessagesByLabelPath(context.Background(), tt.config, tt.labelPath, tt.maxResults)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
				for i, expectedMsg := range tt.expectedResult {
					assert.Equal(t, expectedMsg.ID, result[i].ID)
					assert.Equal(t, expectedMsg.Subject, result[i].Subject)
					assert.Equal(t, expectedMsg.From, result[i].From)
					assert.Equal(t, expectedMsg.To, result[i].To)
					assert.Equal(t, expectedMsg.Body, result[i].Body)
				}
			}

			// モックの期待が満たされたかチェック
			mockGmailAuthService.AssertExpectations(t)
			mockGmailMessageService.AssertExpectations(t)
		})
	}
}

func TestGmailMessageUseCaseImpl_GetMessage(t *testing.T) {
	tests := []struct {
		name           string
		config         domain.GmailAuthConfig
		messageID      string
		setupMock      func(*MockGmailMessageService)
		expectedResult *domain.GmailMessage
		expectedError  string
	}{
		{
			name: "正常系_メッセージ詳細取得成功",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			messageID: "msg1",
			setupMock: func(mockService *MockGmailMessageService) {
				expectedMessage := &domain.GmailMessage{
					ID:      "msg1",
					Subject: "テストメール詳細",
					From:    "sender@example.com",
					To:      []string{"recipient@example.com"},
					Date:    time.Now(),
					Body:    "テストメール詳細の本文",
				}
				mockService.On("GetMessage", mock.Anything, mock.AnythingOfType("domain.GmailCredential"), "test-app", "msg1").Return(expectedMessage, nil)
			},
			expectedResult: &domain.GmailMessage{
				ID:      "msg1",
				Subject: "テストメール詳細",
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Body:    "テストメール詳細の本文",
			},
		},
		{
			name: "異常系_メッセージIDが空",
			config: domain.GmailAuthConfig{
				ClientSecretPath:  "test-client-secret.json",
				CredentialsFolder: "test-credentials",
				ApplicationName:   "test-app",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
				UserID:            "user",
			},
			messageID:     "",
			setupMock:     func(mockService *MockGmailMessageService) {},
			expectedError: "メッセージIDが指定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockGmailAuthService := &MockGmailAuthService{}
			mockGmailMessageService := &MockGmailMessageService{}
			tt.setupMock(mockGmailMessageService)

			// 認証が必要な場合のモック設定
			if tt.expectedError == "" || tt.expectedError != "メッセージIDが指定されていません" {
				credential := domain.GmailCredential{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					TokenType:    "Bearer",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				mockGmailAuthService.On("LoadCredentials", tt.config.CredentialsFolder, "user").Return(&credential, nil)
			}

			useCase := NewGmailMessageUseCase(mockGmailAuthService, mockGmailMessageService)

			// Act
			result, err := useCase.GetMessage(context.Background(), tt.config, tt.messageID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.Subject, result.Subject)
				assert.Equal(t, tt.expectedResult.From, result.From)
				assert.Equal(t, tt.expectedResult.To, result.To)
				assert.Equal(t, tt.expectedResult.Body, result.Body)
			}

			// モックの期待が満たされたかチェック
			mockGmailAuthService.AssertExpectations(t)
			mockGmailMessageService.AssertExpectations(t)
		})
	}
}
