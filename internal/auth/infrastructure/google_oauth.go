// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルはGoogle OAuth2サービスの実装を定義します。
package infrastructure

import (
	"business/internal/auth/application"
	"business/internal/auth/domain"
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// googleOAuthService はGoogle OAuth2サービスの実装です
type googleOAuthService struct {
	config *oauth2.Config
}

// NewGoogleOAuthService は新しいGoogle OAuth2サービスを作成します
func NewGoogleOAuthService(authConfig domain.AuthConfig) application.GoogleOAuthService {
	config := &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  authConfig.RedirectURL,
		Scopes:       authConfig.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &googleOAuthService{
		config: config,
	}
}

// GetAuthURL は認証URLを生成します
func (s *googleOAuthService) GetAuthURL(state string) string {
	// TODO: TDDで実装予定
	return ""
}

// ExchangeCode は認証コードをアクセストークンに交換します
func (s *googleOAuthService) ExchangeCode(ctx context.Context, code string) (*domain.GoogleAuthResponse, error) {
	// TODO: TDDで実装予定
	return nil, fmt.Errorf("not implemented")
}

// GetUserInfo はアクセストークンを使用してユーザー情報を取得します
func (s *googleOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*domain.GoogleUserInfo, error) {
	// TODO: TDDで実装予定
	return nil, fmt.Errorf("not implemented")
}
