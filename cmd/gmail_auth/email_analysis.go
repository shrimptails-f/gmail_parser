// Package main はGmail認証のコマンドラインアプリケーションです。
// このファイルはメール分析機能を提供します。
package main

import (
	emailstoredi "business/internal/emailstore/di"
	"business/internal/gmail/domain"
	aiapp "business/internal/openai/application"
	openaidomain "business/internal/openai/domain"
	aiinfra "business/internal/openai/infrastructure"
	"business/tools/mysql"
	"context"
	"fmt"
	"os"
)

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

// saveEmailAnalysisResult はメール分析結果をDBに保存します
func saveEmailAnalysisResult(ctx context.Context, result *openaidomain.EmailAnalysisResult) error {
	// MySQL接続を作成
	mysqlConn, err := mysql.New()
	if err != nil {
		return fmt.Errorf("MySQL接続エラー: %w", err)
	}

	// EmailStoreUseCaseを作成
	emailStoreUseCase := emailstoredi.ProvideEmailStoreDependencies(mysqlConn.DB)

	// メール分析結果を保存
	if err := emailStoreUseCase.SaveEmailAnalysisResult(ctx, result); err != nil {
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	fmt.Printf("メール分析結果をDBに保存しました: %s\n", result.GmailID)
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

			// 入場時期・開始時期
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

			// 入場時期・開始時期
			if startDates, ok := emailProject["入場時期・開始時期"]; ok {
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
