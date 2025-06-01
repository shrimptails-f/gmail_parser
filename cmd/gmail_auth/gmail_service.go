// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail APIサービステスト機能を提供します。
package main

import (
	"business/internal/gmail/application"
	"business/internal/gmail/infrastructure"
	"business/tools/logger"
	"context"
	"fmt"

	"google.golang.org/api/gmail/v1"
)

// testGmailService はGmail APIサービスを作成してテストします
func testGmailService(ctx context.Context, l *logger.Logger) error {
	// Gmail認証設定を作成
	config, err := createGmailAuthConfig()
	if err != nil {
		return err
	}

	// Gmail認証サービスとユースケースを作成
	gmailAuthService := infrastructure.NewGmailAuthService()
	gmailAuthUseCase := application.NewGmailAuthUseCase(gmailAuthService)

	// Gmail APIサービスを作成
	service, err := gmailAuthUseCase.CreateGmailService(ctx, *config)
	if err != nil {
		return err
	}

	// Gmail APIサービスをテスト
	gmailService, ok := service.(*gmail.Service)
	if !ok {
		return fmt.Errorf("gmail APIサービスの型変換に失敗しました")
	}

	// ユーザープロファイルを取得してテスト
	profile, err := gmailService.Users.GetProfile("me").Do()
	if err != nil {
		return fmt.Errorf("ユーザープロファイルの取得に失敗しました: %w", err)
	}

	// 結果を表示
	fmt.Printf("Gmail APIサービステスト成功!\n")
	fmt.Printf("メールアドレス: %s\n", profile.EmailAddress)
	fmt.Printf("メッセージ総数: %d\n", profile.MessagesTotal)
	fmt.Printf("スレッド総数: %d\n", profile.ThreadsTotal)
	fmt.Printf("履歴ID: %d\n", profile.HistoryId)

	return nil
}
