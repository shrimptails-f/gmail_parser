// Package application は認証機能のアプリケーション層を提供します。
// このファイルは認証機能で使用するインターフェースを定義します。
package application

import (
	"business/internal/gmail/domain"
	"context"
)

// GmailAuthUseCase はGmail認証機能のユースケースインターフェースです
type GmailAuthUseCase interface {
	// AuthenticateGmail はGmail認証を実行します
	AuthenticateGmail(ctx context.Context, config domain.GmailAuthConfig) (*domain.GmailAuthResult, error)

	// CreateGmailService はGmail APIサービスを作成します
	CreateGmailService(ctx context.Context, config domain.GmailAuthConfig) (interface{}, error)
}

// GmailMessageUseCase はGmailメッセージ取得機能のユースケースインターフェースです
type GmailMessageUseCase interface {
	// GetMessages はメッセージ一覧を取得します
	GetMessages(ctx context.Context, config domain.GmailAuthConfig, maxResults int64) ([]domain.GmailMessage, error)

	// GetMessage は指定されたIDのメッセージを取得します
	GetMessage(ctx context.Context, config domain.GmailAuthConfig, messageID string) (*domain.GmailMessage, error)

	// GetMessagesByLabelPath は指定されたラベルパスのメッセージ一覧を取得します
	GetMessagesByLabelPath(ctx context.Context, config domain.GmailAuthConfig, labelPath string, maxResults int64) ([]domain.GmailMessage, error)
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
