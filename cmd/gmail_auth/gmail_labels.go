// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmailラベル取得機能を提供します。
package main

import (
	"business/internal/gmail/infrastructure"
	"business/tools/logger"
	"context"
	"fmt"
)

// testGmailLabels はGmailラベル一覧を取得してテストします
func testGmailLabels(ctx context.Context, l *logger.Logger) error {
	// Gmail認証設定を作成
	config, err := createGmailAuthConfig()
	if err != nil {
		return err
	}

	// サービスとユースケースを作成
	gmailAuthService := infrastructure.NewGmailAuthService()
	gmailMessageService := infrastructure.NewGmailMessageService()

	// 認証情報を取得
	credential, err := gmailAuthService.LoadCredentials(config.CredentialsFolder, config.UserID)
	if err != nil {
		return fmt.Errorf("認証情報の読み込みに失敗しました: %w", err)
	}

	// ラベル一覧を取得
	labels, err := gmailMessageService.GetLabels(ctx, *credential, config.ApplicationName)
	if err != nil {
		return fmt.Errorf("ラベル一覧の取得に失敗しました: %w", err)
	}

	// 結果を表示
	fmt.Printf("Gmailラベル一覧取得テスト成功!\n")
	fmt.Printf("取得したラベル数: %d\n\n", len(labels))

	for i, label := range labels {
		fmt.Printf("=== ラベル %d ===\n", i+1)
		fmt.Printf("ID: %s\n", label.ID)
		fmt.Printf("名前: %s\n", label.Name)
		fmt.Printf("タイプ: %s\n", label.Type)
		fmt.Println()
	}

	return nil
}
