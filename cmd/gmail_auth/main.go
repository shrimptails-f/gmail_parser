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
		label := os.Args[2]
		fmt.Printf("指定ラベル: %s\n", label)

		strSinceDaysAgo := os.Args[3]
		fmt.Printf("日付調整: %s\n", strSinceDaysAgo)
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
		fmt.Printf("%s以降のメールを取得します\n", start.Format("2006-01-02 15:04"))

		if len(os.Args) < 4 {
			fmt.Println("エラー: 何日前から取得するか指定してください")
			fmt.Println("使用例: 前日から取得する場合")
			fmt.Println("go run main.go gmail-messages-by-label 営業/案件 -1")
			fmt.Println("使用例: 当日分を取得する場合")
			fmt.Println("go run main.go gmail-messages-by-label 営業/案件 0")
			return
		}

		var messages []cd.BasicMessage
		err = container.Invoke(func(ga *ga.GmailUseCase) {
			messages, err = ga.GetMessages(ctx, label, sinceDaysAgo)
			if err != nil {
				fmt.Printf("gメール取得処理失敗: %v \n", err)
				return
			}
		})
		if err != nil {
			fmt.Printf("gメール取得処理失敗: %v \n", err)
			return
		}

		// 結果を表示
		fmt.Printf("ラベル指定Gmailメッセージ取得テスト成功!\n")
		fmt.Printf("取得したメッセージ数: %d\n\n", len(messages))

		var exists bool
		for _, message := range messages {
			err = container.Invoke(func(ea *ea.EmailStoreUseCaseImpl) {
				exists, err = ea.CheckGmailIdExists(message.ID)
				if err != nil {
					fmt.Printf("メール存在確認エラー: %v", err)
					return
				}
			})

			// 既に存在する場合はスキップ
			if exists {
				fmt.Printf("メールID %s は既に処理済みです。字句解析をスキップします。\n", message.ID)
				continue
			} else {
				fmt.Printf("メールID %s を解析します。 \n", message.ID)
			}

			// メール本文の分析を実行
			var analysisResults []cd.AnalysisResult
			err = container.Invoke(func(aiapp *aiapp.UseCase) {
				analysisResults, err = aiapp.AnalyzeEmailContent(ctx, message.Body)
				if err != nil {
					fmt.Printf("メール分析エラー: %v", err)
					return
				}
			})

			// 解析結果を保存形式へ詰め替える。
			results := convertToStructs(message, analysisResults)

			// DB保存
			for _, result := range results {
				err = container.Invoke(func(ea *ea.EmailStoreUseCaseImpl) {
					err = ea.SaveEmailAnalysisResult(result)
					if err != nil {
						fmt.Printf("メール保存エラー: %v", err)
						return
					}
				})
			}
		}

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
