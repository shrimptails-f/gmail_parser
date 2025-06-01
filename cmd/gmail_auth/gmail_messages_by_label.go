// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはラベル指定Gmailメッセージ取得機能を提供します。
package main

import (
	"business/internal/gmail/application"
	"business/internal/gmail/infrastructure"
	"business/tools/logger"
	"context"
	"fmt"
)

// getGmailMessagesByLabel はラベル指定でGmailメッセージを取得してテストします
func getGmailMessagesByLabel(ctx context.Context, l *logger.Logger, labelPath string) error {
	// Gmail認証設定を作成
	config, err := createGmailAuthConfig()
	if err != nil {
		return err
	}

	// サービスとユースケースを作成
	gmailAuthService := infrastructure.NewGmailAuthService()
	gmailMessageService := infrastructure.NewGmailMessageService()
	gmailMessageUseCase := application.NewGmailMessageUseCase(gmailAuthService, gmailMessageService)

	// ラベル指定で当日0時以降のメッセージを全件取得
	messages, err := gmailMessageUseCase.GetAllMessagesByLabelPathFromToday(ctx, *config, labelPath, 50)
	if err != nil {
		return fmt.Errorf("ラベル指定メッセージ一覧の取得に失敗しました: %w", err)
	}

	// 結果を表示
	fmt.Printf("ラベル指定Gmailメッセージ取得テスト成功!\n")
	fmt.Printf("指定ラベル: %s\n", labelPath)
	fmt.Printf("取得したメッセージ数: %d\n\n", len(messages))

	for _, message := range messages {
		// メール分析を実行
		if err := analyzeEmailMessage(ctx, &message); err != nil {
			l.Error(fmt.Errorf("メール分析に失敗しました: %w", err))
		}
		fmt.Println()
	}

	return nil
}
