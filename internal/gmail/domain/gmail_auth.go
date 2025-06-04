// Package domain は認証機能のドメイン層を提供します。
// このファイルはGmail認証に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"strings"
	"time"
)

// GmailAuthConfig はGmail認証設定のドメインモデルです
type GmailAuthConfig struct {
	ClientSecretPath  string   // client-secret.jsonファイルのパス
	CredentialsFolder string   // 認証情報を保存するフォルダ
	Scopes            []string // Gmail APIのスコープ
	ApplicationName   string   // アプリケーション名
	LocalServerPort   int      // ローカルサーバーのポート番号
	UserID            string   // 認証するユーザーID（通常は"user"）
}

// GmailCredential はGmail認証の認証情報を表すドメインモデルです
type GmailCredential struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GmailAuthResult はGmail認証結果のドメインモデルです
type GmailAuthResult struct {
	Credential      GmailCredential `json:"credential"`
	ApplicationName string          `json:"application_name"`
	IsNewAuth       bool            `json:"is_new_auth"`
}

// Gmail認証関連のドメインエラー
var (
	ErrClientSecretNotFound    = errors.New("client-secret.jsonファイルが見つかりません")
	ErrInvalidClientSecret     = errors.New("無効なclient-secret.jsonファイルです")
	ErrCredentialsFolderAccess = errors.New("認証情報フォルダにアクセスできません")
	ErrGmailAuthFailed         = errors.New("Gmail認証に失敗しました")
	ErrGmailServiceCreation    = errors.New("Gmailサービスの作成に失敗しました")
	ErrInvalidScope            = errors.New("無効なスコープです")
)

// NewGmailAuthConfig はGmail認証設定を作成します
func NewGmailAuthConfig(clientSecretPath, credentialsFolder, applicationName string) *GmailAuthConfig {
	return &GmailAuthConfig{
		ClientSecretPath:  clientSecretPath,
		CredentialsFolder: credentialsFolder,
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
		ApplicationName: applicationName,
		LocalServerPort: 5555,
		UserID:          "user",
	}
}

// IsValid はGmail認証設定の妥当性をチェックします
func (c *GmailAuthConfig) IsValid() error {
	if c.ClientSecretPath == "" {
		return ErrClientSecretNotFound
	}
	if c.CredentialsFolder == "" {
		return ErrCredentialsFolderAccess
	}
	if c.ApplicationName == "" {
		return errors.New("アプリケーション名が設定されていません")
	}
	if len(c.Scopes) == 0 {
		return ErrInvalidScope
	}
	return nil
}

// IsExpired は認証情報の有効期限をチェックします
func (c *GmailCredential) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// IsValid は認証情報の妥当性をチェックします
func (c *GmailCredential) IsValid() bool {
	return c.AccessToken != "" && !c.IsExpired()
}

// GmailMessage はGmailメッセージのドメインモデルです
type GmailMessage struct {
	ID      string    `json:"id"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	To      []string  `json:"to"`
	Date    time.Time `json:"date"`
	Body    string    `json:"body"`
}

// ExtractSenderName は From フィールドから送信者名を抽出します
func (t GmailMessage) ExtractSenderName() string {
	if idx := strings.Index(t.From, "<"); idx > 0 {
		return strings.TrimSpace(t.From[:idx])
	}
	return t.From
}

// ExtractEmailAddress は From フィールドからメールアドレスを抽出します
func (t GmailMessage) ExtractEmailAddress() string {
	start := strings.Index(t.From, "<")
	end := strings.Index(t.From, ">")
	if start >= 0 && end > start {
		return t.From[start+1 : end]
	}
	return t.From
}

// GmailInfo はGmailラベル情報のドメインモデルです
type GmailInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
