// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能のエントリーポイントを提供します。
package main

import (
	cd "business/internal/common/domain"
	"business/internal/emailstore/application"
	emailstoredi "business/internal/emailstore/di"
	ga "business/internal/gmail/application"
	gi "business/internal/gmail/infrastructure"
	aiapp "business/internal/openAi/application"
	aiinfra "business/internal/openAi/infrastructure"
	"strconv"

	gc "business/tools/gmail"
	gs "business/tools/gmailService"
	"business/tools/logger"
	db "business/tools/mysql"
	oa "business/tools/openai"
	"business/tools/oswrapper"
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

	osw := oswrapper.New()

	credentialsPath := osw.GetEnv("CLIENT_SECRET_PATH")
	tokenPath := "/data/credentials/token_user.json"
	gs := gs.NewClient()
	svc, err := gs.CreateGmailService(ctx, credentialsPath, tokenPath)
	if err != nil {
		fmt.Printf("gメールAPIクライアント生成に失敗しました:%v \n", err)
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
		if err := executeGmailAuth(ctx, gs, credentialsPath, port); err != nil {
			l.Error(fmt.Errorf("gmail認証に失敗しました: %w", err))
			os.Exit(1)
		}
	case "gmail-messages-by-label":
		// ラベル指定でGmailメッセージを取得してテスト
		if len(os.Args) < 3 {
			fmt.Println("エラー: ラベルパスを指定してください")
			fmt.Println("使用例: go run main.go gmail-messages-by-label 営業/案件")
			os.Exit(1)
		}
		label := os.Args[2]
		fmt.Printf("指定ラベル: %s\n", label)

		gc := gc.NewClient(svc)

		// サービスとユースケースを作成
		gi := gi.New(gc)
		ga := ga.New(gi)

		// DB保存機能郡 インスタンス作成
		dbConn, err := db.New()
		if err != nil {
			fmt.Printf("MySQL接続エラー: %v", err)
			return
		}
		es := emailstoredi.ProvideEmailStoreDependencies(dbConn.DB)

		// OpenAi解析機能群 インスタンス作成
		oa := oa.New(osw.GetEnv("OPENAI_API_KEY"))
		aiapp := aiapp.NewUseCase(aiinfra.NewAnalyzer(oa), osw)

		if err := getGmailMessagesByLabel(ctx, l, ga, es, *aiapp, label); err != nil {
			l.Error(fmt.Errorf("ラベル指定gmailメッセージの取得に失敗しました: %w", err))
			os.Exit(1)
		}

	default:
		printUsage()
	}
}

// executeGmailAuth はGmail認証を実行します
func executeGmailAuth(ctx context.Context, gs *gs.Client, credentialsPath string, port int) error {

	result, err := gs.Authenticate(ctx, credentialsPath, port)
	if err != nil {
		return err
	}

	// 認証結果を表示
	fmt.Printf("Gmail認証成功!\n")
	fmt.Printf("アクセストークン: %s...\n", result.AccessToken[:20])
	fmt.Printf("トークンタイプ: %s\n", result.TokenType)
	fmt.Printf("有効期限: %v\n", result.ExpiresIn)

	return nil
}

// getGmailMessagesByLabel はラベル指定でGmailメッセージを取得してテストします
func getGmailMessagesByLabel(ctx context.Context, l *logger.Logger, ga ga.GmailUseCaseInterface, es application.EmailStoreUseCase, aiapp aiapp.UseCase, label string) error {

	messages, err := ga.GetMessages(ctx, label)
	if err != nil {
		fmt.Printf("gメール取得処理失敗: %v", err)
		return err
	}

	// 結果を表示
	fmt.Printf("ラベル指定Gmailメッセージ取得テスト成功!\n")
	fmt.Printf("取得したメッセージ数: %d\n\n", len(messages))

	for _, message := range messages {

		// メールIDの存在確認
		exists, err := es.CheckGmailIdExists(message.ID)
		if err != nil {
			return fmt.Errorf("メール存在確認エラー: %w", err)
		}

		// 既に存在する場合はスキップ
		if exists {
			fmt.Printf("メールID %s は既に処理済みです。字句解析をスキップします。\n", message.ID)
			continue
		}

		// メール本文の分析を実行
		analysisResults, err := aiapp.AnalyzeEmailContent(ctx, message.Body)
		if err != nil {
			return fmt.Errorf("メール分析エラー: %w", err)
		}

		// 解析結果を保存形式へ詰め替える。
		results := convertToStructs(message, analysisResults)

		// DB保存
		for _, result := range results {

			if err := es.SaveEmailAnalysisResult(result); err != nil {
				return fmt.Errorf("複数案件DB保存エラー: %w", err)
			}
		}
	}

	return nil
}

// convertToStructs は引数を結合して保存する形式へ詰め替えます。
func convertToStructs(message cd.BasicMessage, analysisResults []cd.AnalysisResult) []cd.Email {
	var results []cd.Email

	for _, analysisResult := range analysisResults {
		result := cd.Email{
			GmailID:             message.ID,
			ReceivedDate:        message.Date,
			Summary:             analysisResult.ProjectTitle,
			Subject:             message.Subject,
			From:                message.From,
			FromEmail:           message.ExtractEmailAddress(),
			Body:                message.Body,
			Category:            analysisResult.MailCategory,
			ProjectName:         analysisResult.ProjectTitle,
			StartPeriod:         analysisResult.StartPeriod,
			EndPeriod:           analysisResult.EndPeriod,
			WorkLocation:        analysisResult.WorkLocation,
			PriceFrom:           analysisResult.PriceFrom,
			PriceTo:             analysisResult.PriceTo,
			Languages:           analysisResult.Languages,
			Frameworks:          analysisResult.Frameworks,
			Positions:           analysisResult.Positions,
			WorkTypes:           analysisResult.WorkTypes,
			RequiredSkillsMust:  analysisResult.RequiredSkillsMust,
			RequiredSkillsWant:  analysisResult.RequiredSkillsWant,
			RemoteWorkCategory:  analysisResult.RemoteWorkCategory,
			RemoteWorkFrequency: analysisResult.RemoteWorkFrequency,
		}
		results = append(results, result)
	}

	return results
}

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
