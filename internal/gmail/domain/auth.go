// Package domain は認証機能のドメイン層を提供します。
// このファイルは認証に関するドメインモデルとビジネスルールを定義します。
package domain

import (
	"errors"
	"time"
)

// User はユーザーのドメインモデルです
type User struct {
	UserID    uint32    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GoogleAuthRequest はGoogle認証リクエストのドメインモデルです
type GoogleAuthRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// GoogleAuthResponse はGoogle認証レスポンスのドメインモデルです
type GoogleAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// GoogleUserInfo はGoogleから取得するユーザー情報のドメインモデルです
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// AuthResult は認証結果のドメインモデルです
type AuthResult struct {
	User      User   `json:"user"`
	JWTToken  string `json:"jwt_token"`
	IsNewUser bool   `json:"is_new_user"`
}

// AuthConfig はGoogle OAuth設定のドメインモデルです
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// ドメインエラー
var (
	ErrInvalidAuthCode   = errors.New("無効な認証コードです")
	ErrUserNotFound      = errors.New("ユーザーが見つかりません")
	ErrInvalidToken      = errors.New("無効なトークンです")
	ErrTokenExpired      = errors.New("トークンの有効期限が切れています")
	ErrEmailNotVerified  = errors.New("メールアドレスが認証されていません")
	ErrInvalidGoogleUser = errors.New("無効なGoogleユーザー情報です")
)

// IsValidEmail はメールアドレスの形式をチェックします
func (u *User) IsValidEmail() bool {
	return u.Email != "" && len(u.Email) > 0
}

// GetFullName はフルネームを取得します
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidGoogleUserInfo はGoogleユーザー情報の妥当性をチェックします
func (g *GoogleUserInfo) IsValidGoogleUserInfo() error {
	if g.Email == "" {
		return ErrInvalidGoogleUser
	}
	if !g.VerifiedEmail {
		return ErrEmailNotVerified
	}
	return nil
}

// ToUser はGoogleUserInfoからUserドメインモデルに変換します
func (g *GoogleUserInfo) ToUser() User {
	return User{
		Email:     g.Email,
		FirstName: g.GivenName,
		LastName:  g.FamilyName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
