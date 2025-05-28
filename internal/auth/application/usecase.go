// Package application は認証機能のアプリケーション層を提供します。
// このファイルは認証機能のユースケースを実装します。
package application

import (
	"business/internal/auth/domain"
	"context"
	"fmt"
)

// authUseCase は認証機能のユースケース実装です
type authUseCase struct {
	authRepo      AuthRepository
	googleService GoogleOAuthService
	jwtService    JWTService
}

// NewAuthUseCase は新しい認証ユースケースを作成します
func NewAuthUseCase(
	authRepo AuthRepository,
	googleService GoogleOAuthService,
	jwtService JWTService,
) AuthUseCase {
	return &authUseCase{
		authRepo:      authRepo,
		googleService: googleService,
		jwtService:    jwtService,
	}
}

// GetGoogleAuthURL はGoogle認証URLを取得します
func (uc *authUseCase) GetGoogleAuthURL(state string) string {
	// TODO: TDDで実装予定
	return ""
}

// AuthenticateWithGoogle はGoogleアカウントで認証を行います
func (uc *authUseCase) AuthenticateWithGoogle(ctx context.Context, request domain.GoogleAuthRequest) (*domain.AuthResult, error) {
	// TODO: TDDで実装予定
	return nil, fmt.Errorf("not implemented")
}

// ValidateJWTToken はJWTトークンを検証します
func (uc *authUseCase) ValidateJWTToken(token string) (uint32, error) {
	// TODO: TDDで実装予定
	return 0, fmt.Errorf("not implemented")
}
