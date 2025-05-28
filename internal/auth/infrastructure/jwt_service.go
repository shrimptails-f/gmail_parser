// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルはJWTサービスの実装を定義します。
package infrastructure

import (
	"business/internal/auth/application"
	"fmt"
	"os"
)

// jwtService はJWTサービスの実装です
type jwtService struct {
	secretKey string
	appName   string
}

// NewJWTService は新しいJWTサービスを作成します
func NewJWTService() application.JWTService {
	return &jwtService{
		secretKey: os.Getenv("JWT_SECRET_KEY"),
		appName:   os.Getenv("APP_NAME"),
	}
}

// GenerateToken はユーザーIDからJWTトークンを生成します
func (s *jwtService) GenerateToken(userID uint32) (string, error) {
	// TODO: TDDで実装予定
	return "", fmt.Errorf("not implemented")
}

// ValidateToken はJWTトークンを検証してユーザーIDを取得します
func (s *jwtService) ValidateToken(tokenString string) (uint32, error) {
	// TODO: TDDで実装予定
	return 0, fmt.Errorf("not implemented")
}
