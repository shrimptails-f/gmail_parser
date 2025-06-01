// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルは共通機能を提供します。
package main

import (
	"business/internal/gmail/domain"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getClientSecretPath はclient-secret.jsonファイルのパスを取得します
func getClientSecretPath() string {
	// 環境変数から取得
	if path := os.Getenv("CLIENT_SECRET_PATH"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// カレントディレクトリから検索
	candidates := []string{
		"client-secret.json",
		"credentials/client-secret.json",
		"../client-secret.json",
		"../../client-secret.json",
	}

	for _, candidate := range candidates {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}

// createGmailAuthConfig はGmail認証設定を作成します
func createGmailAuthConfig() (*domain.GmailAuthConfig, error) {
	clientSecretPath := getClientSecretPath()
	if clientSecretPath == "" {
		return nil, fmt.Errorf("client-secret.jsonファイルが見つかりません。カレントディレクトリまたは環境変数CLIENT_SECRET_PATHで指定してください")
	}

	config := domain.NewGmailAuthConfig(
		clientSecretPath,
		"credentials",
		"gmailai",
	)

	return config, nil
}

// truncateString は文字列を指定された長さで切り詰めます
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// extractSenderName は送信者情報から名前を抽出します
func extractSenderName(from string) string {
	// 簡単な実装："<"より前の部分を名前とする
	if idx := strings.Index(from, "<"); idx > 0 {
		return strings.TrimSpace(from[:idx])
	}
	return from
}

// extractEmailAddress は送信者情報からメールアドレスを抽出します
func extractEmailAddress(from string) string {
	// 簡単な実装："<"と">"の間の部分をメールアドレスとする
	start := strings.Index(from, "<")
	end := strings.Index(from, ">")
	if start >= 0 && end > start {
		return from[start+1 : end]
	}
	return from
}

// printUsage は使用方法を表示します
func printUsage() {
	fmt.Println("Gmail認証コマンドラインアプリケーション")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  go run main.go gmail-auth                    # Gmail認証を実行")
	fmt.Println("  go run main.go gmail-service                 # Gmail APIサービスを作成してテスト")
	fmt.Println("  go run main.go gmail-messages                # Gmailメッセージを取得してテスト")
	fmt.Println("  go run main.go gmail-labels                  # Gmailラベル一覧を取得してテスト")
	fmt.Println("  go run main.go gmail-messages-by-label <ラベル> # 指定ラベルのメッセージを取得")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  go run main.go gmail-messages-by-label 営業/案件")
	fmt.Println("")
	fmt.Println("必要なファイル:")
	fmt.Println("  client-secret.json - Google Cloud ConsoleからダウンロードしたOAuth2認証情報")
	fmt.Println("")
	fmt.Println("環境変数:")
	fmt.Println("  CLIENT_SECRET_PATH - client-secret.jsonファイルのパス（オプション）")
	fmt.Println("")
	fmt.Println("注意:")
	fmt.Println("  - 初回実行時はブラウザで認証が必要です")
	fmt.Println("  - 認証情報は credentials/ フォルダに保存されます")
	fmt.Println("  - Gmail API の読み取り専用スコープを使用します")
}
