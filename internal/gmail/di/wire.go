// Package di は認証機能の依存性注入を提供します。
// このファイルは認証機能の依存関係を設定します。
package di

import (
	"business/internal/gmail/application"
	"business/internal/gmail/domain"
	"business/internal/gmail/infrastructure"

	"gorm.io/gorm"
)

// AuthContainer は認証機能の依存関係を管理するコンテナです
type AuthContainer struct {
	AuthUseCase application.AuthUseCase
}

// NewAuthContainer は新しい認証コンテナを作成します
func NewAuthContainer(db *gorm.DB, authConfig domain.AuthConfig) *AuthContainer {
	// インフラストラクチャ層の依存関係を作成
	authRepo := infrastructure.NewAuthRepository(db)
	googleService := infrastructure.NewGoogleOAuthService(authConfig)
	jwtService := infrastructure.NewJWTService()

	// アプリケーション層の依存関係を作成
	authUseCase := application.NewAuthUseCase(authRepo, googleService, jwtService)

	return &AuthContainer{
		AuthUseCase: authUseCase,
	}
}

// GetDefaultAuthConfig はデフォルトのGoogle OAuth設定を返します
func GetDefaultAuthConfig() domain.AuthConfig {
	return domain.AuthConfig{
		ClientID:     "", // 環境変数から設定する必要があります
		ClientSecret: "", // 環境変数から設定する必要があります
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}
