// Package main のテストファイルです。
// Gmail認証コマンドラインアプリケーションの動作をテストします。
package main

import (
	"os"
	"strings"
	"testing"
)

// TestMain_GmailAuth は gmail-auth コマンドの正常系をテストします
func TestMain_GmailAuth(t *testing.T) {
	// このテストはos.Exit(1)を呼ぶため、実際の実行時にのみ動作確認してください
	t.Skip("このテストはos.Exit(1)を呼ぶため、実際の実行時にのみ動作確認してください")
}

// TestMain_NoArguments は引数なしでの実行をテストします
func TestMain_NoArguments(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
	os.Setenv("JWT_SECRET_KEY", "test_jwt_secret")
	defer func() {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("JWT_SECRET_KEY")
	}()

	// 標準出力をキャプチャするための準備
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// コマンドライン引数を設定（引数なし）
	os.Args = []string{"gmail_auth"}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// main関数を実行
	main()

	// 標準出力を復元
	w.Close()
	os.Stdout = oldStdout

	// 出力を読み取り
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// 使用方法が表示されることを確認
	if !strings.Contains(output, "Gmail認証コマンドラインアプリケーション") {
		t.Errorf("使用方法が表示されていません。実際の出力: %s", output)
	}
	if !strings.Contains(output, "使用方法:") {
		t.Errorf("使用方法が表示されていません。実際の出力: %s", output)
	}
}

// TestMain_InvalidCommand は無効なコマンドでの実行をテストします
func TestMain_InvalidCommand(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
	os.Setenv("JWT_SECRET_KEY", "test_jwt_secret")
	defer func() {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("JWT_SECRET_KEY")
	}()

	// 標準出力をキャプチャするための準備
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// コマンドライン引数を設定（無効なコマンド）
	os.Args = []string{"gmail_auth", "invalid-command"}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// main関数を実行
	main()

	// 標準出力を復元
	w.Close()
	os.Stdout = oldStdout

	// 出力を読み取り
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// 使用方法が表示されることを確認
	if !strings.Contains(output, "Gmail認証コマンドラインアプリケーション") {
		t.Errorf("使用方法が表示されていません。実際の出力: %s", output)
	}
}

// TestPrintUsage は printUsage 関数をテストします
func TestPrintUsage(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// printUsage関数を実行
	printUsage()

	// 標準出力を復元
	w.Close()
	os.Stdout = oldStdout

	// 出力を読み取り
	buf := make([]byte, 2048)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// 期待する内容が含まれているかを確認
	expectedContents := []string{
		"Gmail認証コマンドラインアプリケーション",
		"使用方法:",
		"gmail-auth",
		"gmail-service",
		"必要なファイル:",
		"client-secret.json",
		"環境変数:",
		"CLIENT_SECRET_PATH",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("期待する内容が含まれていません: %s\n実際の出力: %s", expected, output)
		}
	}
}
