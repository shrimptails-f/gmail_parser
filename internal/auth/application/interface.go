// Package application は認証機能のアプリケーション層を提供します。
// このファイルは認証機能で使用するインターフェースを定義します。
package application

import (
	"business/internal/auth/domain"
	"context"
)

// AuthRepository は認証機能のリポジトリインターフェースです
type AuthRepository interface {
	// GetUserByEmail はメールアドレスでユーザーを取得します
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)

	// CreateUser は新しいユーザーを作成します
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

// GoogleOAuthService はGoogle OAuth認証サービスのインターフェースです
type GoogleOAuthService interface {
	// GetAuthURL は認証URLを生成します
	GetAuthURL(state string) string

	// ExchangeCode は認証コードをアクセストークンに交換します
	ExchangeCode(ctx context.Context, code string) (*domain.GoogleAuthResponse, error)

	// GetUserInfo はアクセストークンを使用してユーザー情報を取得します
	GetUserInfo(ctx context.Context, accessToken string) (*domain.GoogleUserInfo, error)
}

// JWTService はJWTトークン管理サービスのインターフェースです
type JWTService interface {
	// GenerateToken はユーザーIDからJWTトークンを生成します
	GenerateToken(userID uint32) (string, error)

	// ValidateToken はJWTトークンを検証してユーザーIDを取得します
	ValidateToken(token string) (uint32, error)
}

// GmailAuthService はGmail認証サービスのインターフェースです
type GmailAuthService interface {
	// Authenticate はGmail認証を実行します
	Authenticate(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error)

	// CreateGmailService はGmail APIサービスを作成します
	CreateGmailService(ctx context.Context, credential domain.GmailCredential, applicationName string) (interface{}, error)

	// LoadCredentials は保存された認証情報を読み込みます
	LoadCredentials(credentialsFolder, userID string) (*domain.GmailCredential, error)

	// SaveCredentials は認証情報を保存します
	SaveCredentials(credentialsFolder, userID string, credential domain.GmailCredential) error
}

// GmailAuthUseCase はGmail認証機能のユースケースインターフェースです
type GmailAuthUseCase interface {
	// AuthenticateGmail はGmail認証を実行します
	AuthenticateGmail(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error)

	// CreateGmailService はGmail APIサービスを作成します
	CreateGmailService(ctx context.Context, config domain.GmailAuthConfig) (interface{}, error)
}

// AuthUseCase は認証機能のユースケースインターフェースです
type AuthUseCase interface {
	// GetGoogleAuthURL はGoogle認証URLを取得します
	GetGoogleAuthURL(state string) string

	// AuthenticateWithGoogle はGoogleアカウントで認証を行います
	AuthenticateWithGoogle(ctx context.Context, request domain.GoogleAuthRequest) (*domain.AuthResult, error)

	// ValidateJWTToken はJWTトークンを検証します
	ValidateJWTToken(token string) (uint32, error)
}
