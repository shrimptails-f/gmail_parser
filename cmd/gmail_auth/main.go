// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能のエントリーポイントを提供します。
package main

import (
	"business/internal/auth/application"
	"business/internal/auth/domain"
	"business/internal/auth/infrastructure"
	"business/tools/logger"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/api/gmail/v1"
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
			l.Error(fmt.Errorf("Gmail認証に失敗しました: %w", err))
			os.Exit(1)
		}

	case "gmail-service":
		// Gmail APIサービスを作成してテスト
		if err := testGmailService(ctx, l); err != nil {
			l.Error(fmt.Errorf("Gmail APIサービスのテストに失敗しました: %w", err))
			os.Exit(1)
		}

	default:
		printUsage()
	}
}

// executeGmailAuth はGmail認証を実行します
func executeGmailAuth(ctx context.Context) error {
	// client-secret.jsonファイルのパスを取得
	clientSecretPath := getClientSecretPath()
	if clientSecretPath == "" {
		return fmt.Errorf("client-secret.jsonファイルが見つかりません。カレントディレクトリまたは環境変数CLIENT_SECRET_PATHで指定してください")
	}

	// Gmail認証設定を作成
	config := domain.NewGmailAuthConfig(
		clientSecretPath,
		"credentials",
		"gmailai",
	)

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

// testGmailService はGmail APIサービスを作成してテストします
func testGmailService(ctx context.Context, l *logger.Logger) error {
	// client-secret.jsonファイルのパスを取得
	clientSecretPath := getClientSecretPath()
	if clientSecretPath == "" {
		return fmt.Errorf("client-secret.jsonファイルが見つかりません。カレントディレクトリまたは環境変数CLIENT_SECRET_PATHで指定してください")
	}

	// Gmail認証設定を作成
	config := domain.NewGmailAuthConfig(
		clientSecretPath,
		"credentials",
		"gmailai",
	)

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
		return fmt.Errorf("Gmail APIサービスの型変換に失敗しました")
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

func printUsage() {
	fmt.Println("Gmail認証コマンドラインアプリケーション")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  go run main.go gmail-auth     # Gmail認証を実行")
	fmt.Println("  go run main.go gmail-service  # Gmail APIサービスを作成してテスト")
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
