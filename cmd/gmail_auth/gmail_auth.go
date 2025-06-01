// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能を提供します。
package main

import (
	"business/internal/gmail/application"
	"business/internal/gmail/infrastructure"
	"context"
	"fmt"
)

// executeGmailAuth はGmail認証を実行します
func executeGmailAuth(ctx context.Context) error {
	// Gmail認証設定を作成
	config, err := createGmailAuthConfig()
	if err != nil {
		return err
	}

	// Gmail認証サービスとユースケースを作成
	gmailAuthService := infrastructure.NewGmailAuthService()
	gmailAuthUseCase := application.NewGmailAuthUseCase(gmailAuthService)

	// Gmail認証を実行
	result, err := gmailAuthUseCase.AuthenticateGmail(ctx, *config)
	if err != nil {
		return err
	}

	// 認証結果を表示
	fmt.Printf("Gmail認証成功!\n")
	fmt.Printf("アプリケーション名: %s\n", result.ApplicationName)
	fmt.Printf("新規認証: %t\n", result.IsNewAuth)
	fmt.Printf("アクセストークン: %s...\n", result.Credential.AccessToken[:20])
	fmt.Printf("トークンタイプ: %s\n", result.Credential.TokenType)
	fmt.Printf("有効期限: %s\n", result.Credential.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}
