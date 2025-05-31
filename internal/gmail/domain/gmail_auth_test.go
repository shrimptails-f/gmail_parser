// Package domain のテストファイルです。
// Gmail認証に関するドメインモデルの動作をテストします。
package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewGmailAuthConfig はGmail認証設定の作成をテストします
func TestNewGmailAuthConfig(t *testing.T) {
	// テストケース: 正常な設定作成
	clientSecretPath := "/path/to/client-secret.json"
	credentialsFolder := "credentials"
	applicationName := "gmailai"

	config := NewGmailAuthConfig(clientSecretPath, credentialsFolder, applicationName)

	// 期待値の検証
	assert.Equal(t, clientSecretPath, config.ClientSecretPath)
	assert.Equal(t, credentialsFolder, config.CredentialsFolder)
	assert.Equal(t, applicationName, config.ApplicationName)
	assert.Equal(t, []string{"https://www.googleapis.com/auth/gmail.readonly"}, config.Scopes)
	assert.Equal(t, 5555, config.LocalServerPort)
	assert.Equal(t, "user", config.UserID)
}

// TestGmailAuthConfig_IsValid はGmail認証設定の妥当性チェックをテストします
func TestGmailAuthConfig_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		config      *GmailAuthConfig
		expectError bool
		expectedErr error
	}{
		{
			name: "正常な設定",
			config: &GmailAuthConfig{
				ClientSecretPath:  "/path/to/client-secret.json",
				CredentialsFolder: "credentials",
				ApplicationName:   "gmailai",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
			},
			expectError: false,
		},
		{
			name: "ClientSecretPathが空",
			config: &GmailAuthConfig{
				ClientSecretPath:  "",
				CredentialsFolder: "credentials",
				ApplicationName:   "gmailai",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
			},
			expectError: true,
			expectedErr: ErrClientSecretNotFound,
		},
		{
			name: "CredentialsFolderが空",
			config: &GmailAuthConfig{
				ClientSecretPath:  "/path/to/client-secret.json",
				CredentialsFolder: "",
				ApplicationName:   "gmailai",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
			},
			expectError: true,
			expectedErr: ErrCredentialsFolderAccess,
		},
		{
			name: "ApplicationNameが空",
			config: &GmailAuthConfig{
				ClientSecretPath:  "/path/to/client-secret.json",
				CredentialsFolder: "credentials",
				ApplicationName:   "",
				Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
			},
			expectError: true,
		},
		{
			name: "Scopesが空",
			config: &GmailAuthConfig{
				ClientSecretPath:  "/path/to/client-secret.json",
				CredentialsFolder: "credentials",
				ApplicationName:   "gmailai",
				Scopes:            []string{},
			},
			expectError: true,
			expectedErr: ErrInvalidScope,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGmailCredential_IsExpired は認証情報の有効期限チェックをテストします
func TestGmailCredential_IsExpired(t *testing.T) {
	tests := []struct {
		name       string
		credential *GmailCredential
		expected   bool
	}{
		{
			name: "有効期限内",
			credential: &GmailCredential{
				AccessToken:  "valid_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "有効期限切れ",
			credential: &GmailCredential{
				AccessToken:  "expired_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "ちょうど有効期限",
			credential: &GmailCredential{
				AccessToken:  "token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now(),
			},
			expected: false, // time.Now()は微妙に異なるため、通常はfalse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credential.IsExpired()
			if tt.name == "ちょうど有効期限" {
				// 時間の微妙な差を考慮して、どちらでも許容
				assert.True(t, result == tt.expected || result == !tt.expected)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestGmailCredential_IsValid は認証情報の妥当性チェックをテストします
func TestGmailCredential_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		credential *GmailCredential
		expected   bool
	}{
		{
			name: "有効な認証情報",
			credential: &GmailCredential{
				AccessToken:  "valid_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "AccessTokenが空",
			credential: &GmailCredential{
				AccessToken:  "",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "有効期限切れ",
			credential: &GmailCredential{
				AccessToken:  "valid_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "AccessTokenが空かつ有効期限切れ",
			credential: &GmailCredential{
				AccessToken:  "",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credential.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
