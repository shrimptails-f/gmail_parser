// Package infrastructure のテストファイルです。
// Gmail認証サービスの動作をテストします。
package infrastructure

import (
	"business/internal/gmail/domain"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// TestNewGmailAuthService はGmail認証サービスの作成をテストします
func TestNewGmailAuthService(t *testing.T) {
	service := NewGmailAuthService()
	assert.NotNil(t, service)
}

// TestGmailAuthService_SaveAndLoadCredentials は認証情報の保存と読み込みをテストします
func TestGmailAuthService_SaveAndLoadCredentials(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テストデータ
	userID := "test_user"
	credential := domain.GmailCredential{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	// 認証情報を保存
	err = service.SaveCredentials(tempDir, userID, credential)
	assert.NoError(t, err)

	// ファイルが作成されたことを確認
	tokenPath := filepath.Join(tempDir, "token_test_user.json")
	assert.FileExists(t, tokenPath)

	// 認証情報を読み込み
	loadedCredential, err := service.LoadCredentials(tempDir, userID)
	assert.NoError(t, err)
	assert.NotNil(t, loadedCredential)

	// 内容が一致することを確認
	assert.Equal(t, credential.AccessToken, loadedCredential.AccessToken)
	assert.Equal(t, credential.RefreshToken, loadedCredential.RefreshToken)
	assert.Equal(t, credential.TokenType, loadedCredential.TokenType)
	// 時間の比較は秒単位で行う（JSONシリアライゼーションの精度の問題）
	assert.WithinDuration(t, credential.ExpiresAt, loadedCredential.ExpiresAt, time.Second)
}

// TestGmailAuthService_LoadCredentials_FileNotExists は存在しないファイルの読み込みをテストします
func TestGmailAuthService_LoadCredentials_FileNotExists(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 存在しないファイルを読み込み
	credential, err := service.LoadCredentials(tempDir, "nonexistent_user")
	assert.Error(t, err)
	assert.Nil(t, credential)
	assert.Equal(t, domain.ErrClientSecretNotFound, err)
}

// TestGmailAuthService_LoadCredentials_InvalidJSON は無効なJSONファイルの読み込みをテストします
func TestGmailAuthService_LoadCredentials_InvalidJSON(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 無効なJSONファイルを作成
	userID := "test_user"
	tokenPath := filepath.Join(tempDir, "token_test_user.json")
	err = ioutil.WriteFile(tokenPath, []byte("invalid json"), 0600)
	assert.NoError(t, err)

	// 無効なJSONファイルを読み込み
	credential, err := service.LoadCredentials(tempDir, userID)
	assert.Error(t, err)
	assert.Nil(t, credential)
	assert.Contains(t, err.Error(), "認証情報のパースに失敗しました")
}

// TestGmailAuthService_SaveCredentials_DirectoryCreation はディレクトリ作成をテストします
func TestGmailAuthService_SaveCredentials_DirectoryCreation(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 存在しないサブディレクトリを指定
	credentialsDir := filepath.Join(tempDir, "credentials", "subfolder")
	userID := "test_user"
	credential := domain.GmailCredential{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	// 認証情報を保存（ディレクトリが自動作成される）
	err = service.SaveCredentials(credentialsDir, userID, credential)
	assert.NoError(t, err)

	// ディレクトリが作成されたことを確認
	assert.DirExists(t, credentialsDir)

	// ファイルが作成されたことを確認
	tokenPath := filepath.Join(credentialsDir, "token_test_user.json")
	assert.FileExists(t, tokenPath)
}

// TestGmailAuthService_CreateGmailService_ValidCredential は有効な認証情報でのGmailサービス作成をテストします
func TestGmailAuthService_CreateGmailService_ValidCredential(t *testing.T) {
	service := NewGmailAuthService()

	// テストデータ（有効な認証情報）
	credential := domain.GmailCredential{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	applicationName := "gmailai"

	// Gmail APIサービスを作成
	// 注意: 実際のAPIコールは行わないため、サービスオブジェクトの作成のみをテスト
	gmailService, err := service.CreateGmailService(ctx, credential, applicationName)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, gmailService)
}

// TestGmailAuthService_Authenticate_ClientSecretNotFound はclient-secret.jsonが存在しない場合をテストします
func TestGmailAuthService_Authenticate_ClientSecretNotFound(t *testing.T) {
	service := NewGmailAuthService()

	// 存在しないclient-secret.jsonを指定
	config := domain.GmailAuthConfig{
		ClientSecretPath:  "/nonexistent/client-secret.json",
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	ctx := context.Background()

	// 認証を実行
	result, err := service.Authenticate(ctx, config)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "client-secret.jsonファイルの読み込みに失敗しました")
}

// TestGmailAuthService_Authenticate_InvalidClientSecret は無効なclient-secret.jsonの場合をテストします
func TestGmailAuthService_Authenticate_InvalidClientSecret(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 無効なclient-secret.jsonファイルを作成
	clientSecretPath := filepath.Join(tempDir, "client-secret.json")
	err = ioutil.WriteFile(clientSecretPath, []byte("invalid json"), 0600)
	assert.NoError(t, err)

	config := domain.GmailAuthConfig{
		ClientSecretPath:  clientSecretPath,
		CredentialsFolder: "credentials",
		ApplicationName:   "gmailai",
		Scopes:            []string{"https://www.googleapis.com/auth/gmail.readonly"},
		LocalServerPort:   5555,
		UserID:            "user",
	}

	ctx := context.Background()

	// 認証を実行
	result, err := service.Authenticate(ctx, config)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "OAuth2設定の作成に失敗しました")
}

// TestGmailAuthService_SaveCredentials_JSONFormat は保存されるJSONフォーマットをテストします
func TestGmailAuthService_SaveCredentials_JSONFormat(t *testing.T) {
	service := NewGmailAuthService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := ioutil.TempDir("", "gmail_auth_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テストデータ
	userID := "test_user"
	expiresAt := time.Now().Add(1 * time.Hour).Truncate(time.Second) // 秒単位で切り捨て
	credential := domain.GmailCredential{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
	}

	// 認証情報を保存
	err = service.SaveCredentials(tempDir, userID, credential)
	assert.NoError(t, err)

	// ファイルの内容を直接読み込んで検証
	tokenPath := filepath.Join(tempDir, "token_test_user.json")
	data, err := ioutil.ReadFile(tokenPath)
	assert.NoError(t, err)

	// JSONをパースして内容を確認
	var savedToken oauth2.Token
	err = json.Unmarshal(data, &savedToken)
	assert.NoError(t, err)

	assert.Equal(t, credential.AccessToken, savedToken.AccessToken)
	assert.Equal(t, credential.RefreshToken, savedToken.RefreshToken)
	assert.Equal(t, credential.TokenType, savedToken.TokenType)
	assert.WithinDuration(t, credential.ExpiresAt, savedToken.Expiry, time.Second)
}
