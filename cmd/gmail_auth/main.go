// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能のエントリーポイントを提供します。
package main

import (
	cd "business/internal/common/domain"
	ea "business/internal/emailstore/application"
	ga "business/internal/gmail/application"
	aiapp "business/internal/openAi/application"
	"strconv"

	"business/internal/di"
	"business/tools/gmail"
	"business/tools/gmailService"
	"business/tools/mysql"
	"business/tools/openai"
	"business/tools/oswrapper"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/dig"
)

func main() {
	// ロガーを初期化
	// l := logger.New("info")

	// コマンドライン引数をチェック
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	osw := oswrapper.New()
	ctx := context.Background()
	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")
	container, err := getDependencies(ctx, osw, credentialsPath)
	if err != nil {
		return
	}

	switch command {
	case "gmail-auth":
		// Gmail認証を実行
		strPort := osw.GetEnv("GMAIL_PORT")
		port, err := strconv.Atoi(strPort)
		if err != nil {
			fmt.Printf("gメールのリダイレクトポートの取得に失敗しました。ENVのGMAIL_PORTを見直してください。: %v \n", err)
			return
		}
		err = container.Invoke(func(gs *gmailService.Client) {
			result, err := gs.Authenticate(ctx, credentialsPath, port)
			if err != nil {
				fmt.Printf("%v \n,", err)
				return
			}
			// 認証結果を表示
			fmt.Printf("Gmail認証成功!\n")
			fmt.Printf("アクセストークン: %s...\n", result.AccessToken[:20])
			fmt.Printf("トークンタイプ: %s\n", result.TokenType)
			fmt.Printf("有効期限: %v\n", result.ExpiresIn)
		})
		if err != nil {
			fmt.Printf("依存性注入に失敗しました。:%v \n", err)
			return
		}

	case "gmail-messages-by-label":
		// ラベル指定でGmailメッセージを取得してテスト
		if len(os.Args) < 3 {
			fmt.Println("エラー: ラベルパスを指定してください")
			fmt.Println("使用例: go run main.go gmail-messages-by-label 営業/案件 0")
			return
		}
		if len(os.Args) < 4 {
			fmt.Println("エラー: 何日前から取得するか指定してください")
			fmt.Println("使用例: 前日から取得する場合")
			fmt.Println("go run main.go gmail-messages-by-label 営業/案件 -1")
			fmt.Println("使用例: 当日分を取得する場合")
			fmt.Println("go run main.go gmail-messages-by-label 営業/案件 0")
			return
		}
		strSinceDaysAgo := os.Args[3]
		sinceDaysAgo, err := strconv.Atoi(strSinceDaysAgo)
		if err != nil {
			fmt.Printf("引数の日付調整値の数値変換に失敗しました。引数を確認してください。: %v \n", err)
			return
		}
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if sinceDaysAgo != 0 {
			start = start.AddDate(0, 0, sinceDaysAgo)
		}

		label := os.Args[2]
		fmt.Printf("指定ラベル: %s\n", label)
		fmt.Printf("日付調整: %s\n", strSinceDaysAgo)
		fmt.Printf("%s以降のメールを取得します\n", start.Format("2006-01-02 15:04"))

		var messages []cd.BasicMessage
		var innerErr error
		err = container.Invoke(func(ga *ga.GmailUseCase) {
			messages, innerErr = ga.GetMessages(ctx, label, sinceDaysAgo)
		})
		if innerErr != nil {
			fmt.Printf("gメール取得処理失敗: %v \n", innerErr)
			return
		}
		if err != nil {
			fmt.Printf("gメール取得処理失敗: %v \n", err)
			return
		}

		fmt.Printf("メール分析を行います。 \n")
		var analysisResults []cd.Email
		var AnalyzeinnerErr error
		err = container.Invoke(func(aiapp *aiapp.UseCase) {
			analysisResults, AnalyzeinnerErr = aiapp.AnalyzeEmailContent(ctx, messages)
		})
		if AnalyzeinnerErr != nil {
			fmt.Printf("メール分析エラー: %v \n", AnalyzeinnerErr)
			return
		}
		if err != nil {
			fmt.Printf("メール分析エラー: %v \n", err)
			return
		}

		fmt.Printf("DBへの保存処理を開始します。")
		for _, email := range analysisResults {
			err = container.Invoke(func(ea *ea.EmailStoreUseCaseImpl) {
				err = ea.SaveEmailAnalysisResult(email)
				if err != nil {
					fmt.Printf("メール保存エラー: %v \n", err)
					return
				}
			})
		}
		fmt.Printf("DBへの保存処理が完了しました。")

	default:
		printUsage()
	}
}

func getDependencies(ctx context.Context, osw *oswrapper.OsWrapper, credentialsPath string) (*dig.Container, error) {
	db, err := mysql.New()
	if err != nil {
		fmt.Printf("DB 初期化時にエラーが発生しました。:%v \n,", err)
		return &dig.Container{}, err
	}
	apiKey := osw.GetEnv("OPENAI_API_KEY")
	oa := openai.New(apiKey)

	gs := gmailService.NewClient()
	tokenPath := "/data/credentials/token_user.json"
	svc, err := gs.CreateGmailService(ctx, credentialsPath, tokenPath)

	gc := gmail.NewClient(svc)

	if err != nil {
		fmt.Printf("gメールAPIクライアント生成に失敗しました:%v \n", err)
		return &dig.Container{}, err
	}

	return di.BuildContainer(db, oa, gs, gc, osw), nil

}

func printUsage() {
	fmt.Println("Gmail認証コマンドラインアプリケーション")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  go run main.go gmail-auth                    # Gmail認証を実行")
	fmt.Println("  go run main.go gmail-messages-by-label <ラベル> <日付調整> # 指定ラベルのメッセージを取得")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  使用例: 前日から取得する場合")
	fmt.Println("    go run main.go gmail-messages-by-label 営業/案件 -1")
	fmt.Println("  使用例: 当日分を取得する場合")
	fmt.Println("    go run main.go gmail-messages-by-label 営業/案件 0")
	fmt.Println("")
	fmt.Println("必要なファイル:")
	fmt.Println("  client-secret.json - Google Cloud ConsoleからダウンロードしたOAuth2認証情報")
	fmt.Println("")
	fmt.Println("環境変数:")
	fmt.Println("  LABEL              - Gメールの取得対象となるラベル")
	fmt.Println("  CLIENT_SECRET_PATH - client-secret.jsonファイルのパス(オプション)")
	fmt.Println("  OPENAI_API_KEY     - openAi API秘密鍵")
	fmt.Println("")
	fmt.Println("注意:")
	fmt.Println("  - 初回実行時はブラウザで認証が必要です")
	fmt.Println("  - 認証情報は /data/credentials/ フォルダに保存されます")
	fmt.Println("  - Gmail API の読み取り専用スコープを使用します")
}
