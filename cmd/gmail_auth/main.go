// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能のエントリーポイントを提供します。
package main

import (
	"business/tools/logger"
	"context"
	"fmt"
	"os"
)

func main() {
	// ロガーを初期化
	l := logger.New("info")

	// コマンドライン引数をチェック
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	ctx := context.Background()

	switch command {
	case "gmail-auth":
		// Gmail認証を実行
		if err := executeGmailAuth(ctx); err != nil {
			l.Error(fmt.Errorf("gmail認証に失敗しました: %w", err))
			os.Exit(1)
		}

	case "gmail-service":
		// Gmail APIサービスを作成してテスト
		if err := testGmailService(ctx, l); err != nil {
			l.Error(fmt.Errorf("gmail APIサービスのテストに失敗しました: %w", err))
			os.Exit(1)
		}

	case "gmail-messages":
		// Gmailメッセージを取得してテスト
		if err := testGmailMessages(ctx, l); err != nil {
			l.Error(fmt.Errorf("gmailメッセージの取得に失敗しました: %w", err))
			os.Exit(1)
		}

	case "gmail-labels":
		// Gmailラベル一覧を取得してテスト
		if err := testGmailLabels(ctx, l); err != nil {
			l.Error(fmt.Errorf("gmailラベル一覧の取得に失敗しました: %w", err))
			os.Exit(1)
		}

	case "gmail-messages-by-label":
		// ラベル指定でGmailメッセージを取得してテスト
		if len(os.Args) < 3 {
			fmt.Println("エラー: ラベルパスを指定してください")
			fmt.Println("使用例: go run main.go gmail-messages-by-label 営業/案件")
			os.Exit(1)
		}
		labelPath := os.Args[2]
		if err := getGmailMessagesByLabel(ctx, l, labelPath); err != nil {
			l.Error(fmt.Errorf("ラベル指定gmailメッセージの取得に失敗しました: %w", err))
			os.Exit(1)
		}

	default:
		printUsage()
	}
}
