// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはGmail認証機能のエントリーポイントを提供します。
package main

import (
	emailstoredi "business/internal/emailstore/di"
	"business/internal/gmail/application"
	"business/internal/gmail/domain"
	"business/internal/gmail/infrastructure"
	aiapp "business/internal/openai/application"
	openaidomain "business/internal/openai/domain"
	aiinfra "business/internal/openai/infrastructure"

	"business/tools/logger"
	"business/tools/mysql"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		return fmt.Errorf("gmail APIサービスの型変換に失敗しました")
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

// testGmailMessages はGmailメッセージを取得してテストします
func testGmailMessages(ctx context.Context, l *logger.Logger) error {
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

// testGmailLabels はGmailラベル一覧を取得してテストします
func testGmailLabels(ctx context.Context, l *logger.Logger) error {
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

// getGmailMessagesByLabel はラベル指定でGmailメッセージを取得してテストします
func getGmailMessagesByLabel(ctx context.Context, l *logger.Logger, labelPath string) error {
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

// analyzeEmailMessage はメールメッセージを分析します
func analyzeEmailMessage(ctx context.Context, message *domain.GmailMessage) error {
	// MySQL接続を作成してメールID存在確認
	mysqlConn, err := mysql.New()
	if err != nil {
		return fmt.Errorf("MySQL接続エラー: %w", err)
	}

	// EmailStoreUseCaseを作成
	emailStoreUseCase := emailstoredi.ProvideEmailStoreDependencies(mysqlConn.DB)

	// メールIDの存在確認
	exists, err := emailStoreUseCase.CheckGmailIdExists(ctx, message.ID)
	if err != nil {
		return fmt.Errorf("メール存在確認エラー: %w", err)
	}

	// 既に存在する場合はスキップ
	if exists {
		fmt.Printf("メールID %s は既に処理済みです。字句解析をスキップします。\n", message.ID)
		return nil
	}

	// サービスを作成
	promptService := aiinfra.NewFilePromptService("prompts")
	// OpenAI APIキーを環境変数から取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY環境変数が設定されていません")
	}
	textAnalysisService := aiinfra.NewOpenAIService(apiKey)
	textAnalysisUseCase := aiapp.NewTextAnalysisUseCase(textAnalysisService, promptService)

	// 複数案件対応のメール分析を実行
	results, err := textAnalysisUseCase.AnalyzeEmailTextMultiple(ctx, message.Body, message.ID, message.Subject)
	if err != nil {
		return fmt.Errorf("複数案件メール分析エラー: %w", err)
	}

	// 複数件対応のDB保存処理を実行
	if err := saveEmailAnalysisMultipleResults(ctx, message, results); err != nil {
		return fmt.Errorf("複数案件DB保存エラー: %w", err)
	}

	return nil
}

// saveEmailAnalysisMultipleResults は複数案件対応のメール分析結果をDBに保存します
func saveEmailAnalysisMultipleResults(ctx context.Context, message *domain.GmailMessage, results []*openaidomain.TextAnalysisResult) error {
	// MySQL接続を作成
	mysqlConn, err := mysql.New()
	if err != nil {
		return fmt.Errorf("MySQL接続エラー: %w", err)
	}

	// EmailStoreUseCaseを作成
	emailStoreUseCase := emailstoredi.ProvideEmailStoreDependencies(mysqlConn.DB)

	// 複数案件対応の結果を作成
	multipleResult := openaidomain.NewEmailAnalysisMultipleResult(
		message.ID,
		message.Subject,
		extractSenderName(message.From),
		extractEmailAddress(message.From),
		message.Body,
		message.Date,
	)

	// メール区分を設定（デフォルトで案件）
	multipleResult.MailCategory = "案件"

	// AI解析結果から案件情報を抽出
	for _, result := range results {
		// 営業案件メール情報がある場合は案件として追加
		if hasProjectInfo(result) {
			project := convertToProjectAnalysisResult(result)
			multipleResult.AddProject(project)
		}
	}

	// 案件情報がない場合はデフォルトの案件を追加
	if !multipleResult.HasProjects() {
		defaultProject := openaidomain.ProjectAnalysisResult{
			ProjectName:         "不明",
			StartPeriod:         []string{},
			EndPeriod:           "",
			WorkLocation:        "",
			PriceFrom:           nil,
			PriceTo:             nil,
			Languages:           []string{},
			Frameworks:          []string{},
			Positions:           []string{},
			WorkTypes:           []string{},
			RequiredSkillsMust:  []string{},
			RequiredSkillsWant:  []string{},
			RemoteWorkCategory:  "",
			RemoteWorkFrequency: nil,
		}
		multipleResult.AddProject(defaultProject)
	}

	// 複数案件対応のメール分析結果を保存
	if err := emailStoreUseCase.SaveEmailAnalysisMultipleResult(ctx, multipleResult); err != nil {
		return fmt.Errorf("複数案件メール保存エラー: %w", err)
	}

	fmt.Printf("複数案件対応メール分析結果をDBに保存しました: %s (案件数: %d件)\n", multipleResult.GmailID, multipleResult.GetProjectCount())
	return nil
}

// hasProjectInfo は案件情報が含まれているかチェックします
func hasProjectInfo(result *openaidomain.TextAnalysisResult) bool {
	// キーワードやエンティティから案件情報を判定
	return len(result.Keywords) > 0 || len(result.Entities) > 0
}

// convertToProjectAnalysisResult はTextAnalysisResultを案件分析結果に変換します
func convertToProjectAnalysisResult(result *openaidomain.TextAnalysisResult) openaidomain.ProjectAnalysisResult {
	project := openaidomain.ProjectAnalysisResult{
		ProjectName:         result.Summary, // 要約を案件名として使用
		StartPeriod:         []string{},
		EndPeriod:           "",
		WorkLocation:        "",
		PriceFrom:           nil,
		PriceTo:             nil,
		Languages:           []string{},
		Frameworks:          []string{},
		Positions:           []string{},
		WorkTypes:           []string{},
		RequiredSkillsMust:  []string{},
		RequiredSkillsWant:  []string{},
		RemoteWorkCategory:  "",
		RemoteWorkFrequency: nil,
	}

	// デバッグ: RawResponseの内容を確認
	fmt.Printf("=== デバッグ: RawResponse ===\n")
	for key, value := range result.RawResponse {
		fmt.Printf("Key: %s, Value: %+v\n", key, value)
	}
	fmt.Printf("========================\n")

	// RawResponseからOpenAIの解析結果を取得
	if rawResponse, exists := result.RawResponse["email_project"]; exists {
		// openai_service.goのEmailProjectResponseを使用して型アサーション
		if emailProject, ok := rawResponse.(aiinfra.EmailProjectResponse); ok {
			fmt.Printf("=== 構造体として解析成功 ===\n")
			fmt.Printf("ProjectName: %s\n", emailProject.ProjectName)
			fmt.Printf("SalaryFrom: %d\n", emailProject.SalaryFrom)
			fmt.Printf("SalaryTo: %d\n", emailProject.SalaryTo)
			fmt.Printf("WorkLocation: %s\n", emailProject.WorkLocation)

			// 案件名
			if emailProject.ProjectName != "" {
				project.ProjectName = emailProject.ProjectName
			}

			// 単価FROM
			if emailProject.SalaryFrom > 0 {
				project.PriceFrom = &emailProject.SalaryFrom
			}

			// 単価TO
			if emailProject.SalaryTo > 0 {
				project.PriceTo = &emailProject.SalaryTo
			}

			// 勤務場所
			if emailProject.WorkLocation != "" {
				project.WorkLocation = emailProject.WorkLocation
			}

			// 終了時期
			if emailProject.EndDate != "" {
				project.EndPeriod = emailProject.EndDate
			}

			// 開始時期
			project.StartPeriod = emailProject.StartDates

			// 言語
			project.Languages = emailProject.Languages

			// フレームワーク
			project.Frameworks = emailProject.Frameworks

			// ポジション
			project.Positions = emailProject.Positions

			// 求めるスキル MUST
			project.RequiredSkillsMust = emailProject.RequiredSkills

			// 求めるスキル WANT
			project.RequiredSkillsWant = emailProject.PreferredSkills

			// リモートワーク区分
			if emailProject.RemoteWorkType != "" {
				project.RemoteWorkCategory = emailProject.RemoteWorkType
			}

			// リモートワークの頻度
			if emailProject.RemoteFrequency != "" {
				project.RemoteWorkFrequency = &emailProject.RemoteFrequency
			}
		} else if emailProject, ok := rawResponse.(map[string]interface{}); ok {
			// 案件名
			if projectName, ok := emailProject["案件名"].(string); ok && projectName != "" {
				project.ProjectName = projectName
			}

			// 単価FROM
			if salaryFrom, ok := emailProject["単価FROM"]; ok {
				switch v := salaryFrom.(type) {
				case float64:
					if v > 0 {
						intValue := int(v)
						project.PriceFrom = &intValue
					}
				case int:
					if v > 0 {
						project.PriceFrom = &v
					}
				}
			}

			// 単価TO
			if salaryTo, ok := emailProject["単価TO"]; ok {
				switch v := salaryTo.(type) {
				case float64:
					if v > 0 {
						intValue := int(v)
						project.PriceTo = &intValue
					}
				case int:
					if v > 0 {
						project.PriceTo = &v
					}
				}
			}

			// 勤務場所
			if workLocation, ok := emailProject["勤務場所"].(string); ok && workLocation != "" {
				project.WorkLocation = workLocation
			}

			// 終了時期
			if endDate, ok := emailProject["終了時期"].(string); ok && endDate != "" {
				project.EndPeriod = endDate
			}

			// 開始時期
			if startDates, ok := emailProject["開始時期"]; ok {
				if startDatesArray, ok := startDates.([]interface{}); ok {
					for _, startDate := range startDatesArray {
						if startDateStr, ok := startDate.(string); ok && startDateStr != "" {
							project.StartPeriod = append(project.StartPeriod, startDateStr)
						}
					}
				}
			}

			// 言語
			if languages, ok := emailProject["言語"]; ok {
				if languagesArray, ok := languages.([]interface{}); ok {
					for _, lang := range languagesArray {
						if langStr, ok := lang.(string); ok && langStr != "" {
							project.Languages = append(project.Languages, langStr)
						}
					}
				}
			}

			// フレームワーク
			if frameworks, ok := emailProject["フレームワーク"]; ok {
				if frameworksArray, ok := frameworks.([]interface{}); ok {
					for _, fw := range frameworksArray {
						if fwStr, ok := fw.(string); ok && fwStr != "" {
							project.Frameworks = append(project.Frameworks, fwStr)
						}
					}
				}
			}

			// ポジション
			if positions, ok := emailProject["ポジション"]; ok {
				if positionsArray, ok := positions.([]interface{}); ok {
					for _, pos := range positionsArray {
						if posStr, ok := pos.(string); ok && posStr != "" {
							project.Positions = append(project.Positions, posStr)
						}
					}
				}
			}

			// 求めるスキル MUST
			if mustSkills, ok := emailProject["求めるスキル MUST"]; ok {
				if mustSkillsArray, ok := mustSkills.([]interface{}); ok {
					for _, skill := range mustSkillsArray {
						if skillStr, ok := skill.(string); ok && skillStr != "" {
							project.RequiredSkillsMust = append(project.RequiredSkillsMust, skillStr)
						}
					}
				}
			}

			// 求めるスキル WANT
			if wantSkills, ok := emailProject["求めるスキル WANT"]; ok {
				if wantSkillsArray, ok := wantSkills.([]interface{}); ok {
					for _, skill := range wantSkillsArray {
						if skillStr, ok := skill.(string); ok && skillStr != "" {
							project.RequiredSkillsWant = append(project.RequiredSkillsWant, skillStr)
						}
					}
				}
			}

			// リモートワーク区分
			if remoteType, ok := emailProject["リモートワーク区分"].(string); ok && remoteType != "" {
				project.RemoteWorkCategory = remoteType
			}

			// リモートワークの頻度
			if remoteFreq, ok := emailProject["リモートワークの頻度"].(string); ok && remoteFreq != "" {
				project.RemoteWorkFrequency = &remoteFreq
			}
		}
	}

	// フォールバック: キーワードから技術情報を抽出（OpenAIの解析結果がない場合）
	if len(project.Languages) == 0 && len(project.Frameworks) == 0 {
		for _, keyword := range result.Keywords {
			switch keyword.Category {
			case "言語":
				project.Languages = append(project.Languages, keyword.Text)
			case "フレームワーク":
				project.Frameworks = append(project.Frameworks, keyword.Text)
			}
		}
	}

	// フォールバック: エンティティからポジション情報を抽出（OpenAIの解析結果がない場合）
	if len(project.Positions) == 0 {
		for _, entity := range result.Entities {
			if entity.Type == "POSITION" {
				project.Positions = append(project.Positions, entity.Name)
			}
		}
	}

	return project
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
