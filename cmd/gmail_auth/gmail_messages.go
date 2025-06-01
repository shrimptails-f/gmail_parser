// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmailメッセージ取得機能を提供します。
package main

import (
	"business/internal/gmail/application"
	"business/internal/gmail/infrastructure"
	"business/tools/logger"
	"context"
	"fmt"
)

// testGmailMessages はGmailメッセージを取得してテストします
func testGmailMessages(ctx context.Context, l *logger.Logger) error {
	// Gmail認証設定を作成
	config, err := createGmailAuthConfig()
	if err != nil {
		return err
	}

	// サービスとユースケースを作成
	gmailAuthService := infrastructure.NewGmailAuthService()
	gmailMessageService := infrastructure.NewGmailMessageService()
	gmailMessageUseCase := application.NewGmailMessageUseCase(gmailAuthService, gmailMessageService)

	// メッセージ一覧を取得（最大5件）
	messages, err := gmailMessageUseCase.GetMessages(ctx, *config, 5)
	if err != nil {
		return fmt.Errorf("メッセージ一覧の取得に失敗しました: %w", err)
	}

	// 結果を表示
	fmt.Printf("Gmailメッセージ取得テスト成功!\n")
	fmt.Printf("取得したメッセージ数: %d\n\n", len(messages))

	return nil
}
