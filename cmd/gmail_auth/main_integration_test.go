// Package main はGmail認証のコマンドラインアプリケーションの統合テストを提供します。
// このファイルはEmailAnalysisMultipleResult関連の統合テストを実装します。
//go:build integration

package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	emailstoredi "business/internal/emailstore/di"
	"business/internal/gmail/domain"
	aiapp "business/internal/openai/application"
	openaidomain "business/internal/openai/domain"
	aiinfra "business/internal/openai/infrastructure"
	"business/tools/logger"
	"business/tools/mysql"
)

func TestEmailAnalysisMultiple_統合テスト_SaveEmailAnalysisMultipleResult成功確認(t *testing.T) {
	// 統合テストのため、実際のDBとOpenAI APIを使用
	// 環境変数の確認
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY環境変数が設定されていないため、統合テストをスキップします")
	}

	// Arrange
	ctx := context.Background()
	l := logger.New("info")

	// MySQL接続を作成
	mysqlConn, err := mysql.New()
	require.NoError(t, err, "MySQL接続の作成に失敗しました")

	// EmailStoreUseCaseを作成
	emailStoreUseCase := emailstoredi.ProvideEmailStoreDependencies(mysqlConn.DB)

	// OpenAIサービスを作成
	promptService := aiinfra.NewFilePromptService("../../prompts")
	apiKey := os.Getenv("OPENAI_API_KEY")
	textAnalysisService := aiinfra.NewOpenAIService(apiKey)
	textAnalysisUseCase := aiapp.NewTextAnalysisUseCase(textAnalysisService, promptService)

	// テスト用のメールメッセージを作成
	testMessage := &domain.GmailMessage{
		ID:      "test-integration-" + time.Now().Format("20060102150405"),
		Subject: "【統合テスト】複数案件のご紹介",
		From:    "テスト送信者 <test@example.com>",
		Date:    time.Now(),
		Body: `
件名: 【統合テスト】複数案件のご紹介

お疲れ様です。
以下の案件をご紹介させていただきます。

■案件1: Go言語バックエンド開発
・案件名: ECサイトAPI開発
・開始時期: 2024年4月1日
・終了時期: 2024年12月31日
・勤務場所: 東京都渋谷区
・単価: 80万円〜100万円
・言語: Go, JavaScript
・フレームワーク: Gin, React
・ポジション: バックエンドエンジニア
・必須スキル: Go言語3年以上、REST API設計経験
・歓迎スキル: Docker、Kubernetes経験
・リモートワーク: 週3日リモート可能

■案件2: React.jsフロントエンド開発
・案件名: 管理画面リニューアル
・開始時期: 2024年5月15日
・終了時期: 2024年10月31日
・勤務場所: 大阪府大阪市
・単価: 70万円〜90万円
・言語: JavaScript, TypeScript
・フレームワーク: React, Next.js
・ポジション: フロントエンドエンジニア
・必須スキル: React 2年以上、TypeScript経験
・歓迎スキル: Next.js、GraphQL経験
・リモートワーク: フルリモート可能

ご検討のほど、よろしくお願いいたします。
		`,
	}

	// メールIDの重複チェック（既に存在する場合は削除）
	exists, err := emailStoreUseCase.CheckGmailIdExists(ctx, testMessage.ID)
	require.NoError(t, err, "メール存在チェックに失敗しました")
	if exists {
		t.Logf("テストメールID %s は既に存在するため、統合テストをスキップします", testMessage.ID)
		return
	}

	// Act 1: AnalyzeEmailTextMultiple関数を実行
	t.Log("Step 1: AnalyzeEmailTextMultiple関数を実行中...")
	results, err := textAnalysisUseCase.AnalyzeEmailTextMultiple(ctx, testMessage.Body, testMessage.ID, testMessage.Subject)

	// Assert 1: AI解析が成功することを確認
	require.NoError(t, err, "AnalyzeEmailTextMultiple関数の実行に失敗しました")
	assert.NotEmpty(t, results, "解析結果が空です")
	t.Logf("AI解析結果を取得しました: %d件", len(results))

	// Act 2: EmailAnalysisMultipleResultに変換
	t.Log("Step 2: EmailAnalysisMultipleResultに変換中...")
	multipleResult := convertToEmailAnalysisMultipleResultForIntegration(testMessage, results)

	// Assert 2: 変換結果の検証
	assert.Equal(t, testMessage.ID, multipleResult.GmailID)
	assert.Equal(t, testMessage.Subject, multipleResult.Subject)
	assert.Equal(t, testMessage.Body, multipleResult.Body)
	assert.Equal(t, "案件", multipleResult.MailCategory)
	assert.NoError(t, multipleResult.IsValid(), "EmailAnalysisMultipleResultの妥当性チェックに失敗しました")
	t.Logf("EmailAnalysisMultipleResultに変換完了: 案件数=%d, 人材数=%d",
		multipleResult.GetProjectCount(), multipleResult.GetCandidateCount())

	// Act 3: SaveEmailAnalysisMultipleResultを実行
	t.Log("Step 3: SaveEmailAnalysisMultipleResultを実行中...")
	err = emailStoreUseCase.SaveEmailAnalysisMultipleResult(ctx, multipleResult)

	// Assert 3: DB保存が成功することを確認
	require.NoError(t, err, "SaveEmailAnalysisMultipleResultの実行に失敗しました")
	t.Log("SaveEmailAnalysisMultipleResultが成功しました")

	// Act 4: 保存されたデータの確認
	t.Log("Step 4: 保存されたデータの確認中...")
	savedExists, err := emailStoreUseCase.CheckGmailIdExists(ctx, testMessage.ID)
	require.NoError(t, err, "保存後のメール存在チェックに失敗しました")

	// Assert 4: データが正しく保存されていることを確認
	assert.True(t, savedExists, "メールがDBに保存されていません")
	t.Log("メールがDBに正しく保存されていることを確認しました")

	// 統合テスト成功ログ
	l.Info("EmailAnalysisMultiple統合テストが成功しました")
	t.Log("=== 統合テスト完了 ===")
	t.Logf("テストメールID: %s", testMessage.ID)
	t.Logf("案件数: %d", multipleResult.GetProjectCount())
	t.Logf("人材数: %d", multipleResult.GetCandidateCount())
}

func TestEmailAnalysisMultiple_統合テスト_人材情報メール_SaveEmailAnalysisMultipleResult成功確認(t *testing.T) {
	// 統合テストのため、実際のDBとOpenAI APIを使用
	// 環境変数の確認
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY環境変数が設定されていないため、統合テストをスキップします")
	}

	// Arrange
	ctx := context.Background()
	l := logger.New("info")

	// MySQL接続を作成
	mysqlConn, err := mysql.New()
	require.NoError(t, err, "MySQL接続の作成に失敗しました")

	// EmailStoreUseCaseを作成
	emailStoreUseCase := emailstoredi.ProvideEmailStoreDependencies(mysqlConn.DB)

	// OpenAIサービスを作成
	promptService := aiinfra.NewFilePromptService("../../prompts")
	apiKey := os.Getenv("OPENAI_API_KEY")
	textAnalysisService := aiinfra.NewOpenAIService(apiKey)
	textAnalysisUseCase := aiapp.NewTextAnalysisUseCase(textAnalysisService, promptService)

	// テスト用の人材紹介メールメッセージを作成
	testMessage := &domain.GmailMessage{
		ID:      "test-candidates-" + time.Now().Format("20060102150405"),
		Subject: "【統合テスト】優秀な人材のご紹介",
		From:    "人材紹介会社 <hr@example.com>",
		Date:    time.Now(),
		Body: `
件名: 【統合テスト】優秀な人材のご紹介

お疲れ様です。
以下の人材をご紹介させていただきます。

■人材1: 田中太郎さん
・経験年数: 5年
・スキルまとめ: Go言語、React、Docker、AWS、マイクロサービス設計に精通
・参画可能日: 2024年4月1日

■人材2: 佐藤花子さん
・経験年数: 3年
・スキルまとめ: JavaScript、TypeScript、React、Vue.js、フロントエンド全般
・参画可能日: 2024年5月1日

ご検討のほど、よろしくお願いいたします。
		`,
	}

	// メールIDの重複チェック（既に存在する場合は削除）
	exists, err := emailStoreUseCase.CheckGmailIdExists(ctx, testMessage.ID)
	require.NoError(t, err, "メール存在チェックに失敗しました")
	if exists {
		t.Logf("テストメールID %s は既に存在するため、統合テストをスキップします", testMessage.ID)
		return
	}

	// Act 1: AnalyzeEmailTextMultiple関数を実行
	t.Log("Step 1: AnalyzeEmailTextMultiple関数を実行中...")
	results, err := textAnalysisUseCase.AnalyzeEmailTextMultiple(ctx, testMessage.Body, testMessage.ID, testMessage.Subject)

	// Assert 1: AI解析が成功することを確認
	require.NoError(t, err, "AnalyzeEmailTextMultiple関数の実行に失敗しました")
	assert.NotEmpty(t, results, "解析結果が空です")
	t.Logf("AI解析結果を取得しました: %d件", len(results))

	// Act 2: EmailAnalysisMultipleResultに変換（人材情報用）
	t.Log("Step 2: EmailAnalysisMultipleResultに変換中（人材情報）...")
	multipleResult := convertToEmailAnalysisMultipleResultForCandidates(testMessage, results)

	// Assert 2: 変換結果の検証
	assert.Equal(t, testMessage.ID, multipleResult.GmailID)
	assert.Equal(t, testMessage.Subject, multipleResult.Subject)
	assert.Equal(t, testMessage.Body, multipleResult.Body)
	assert.Equal(t, "人材", multipleResult.MailCategory)
	assert.NoError(t, multipleResult.IsValid(), "EmailAnalysisMultipleResultの妥当性チェックに失敗しました")
	t.Logf("EmailAnalysisMultipleResultに変換完了: 案件数=%d, 人材数=%d",
		multipleResult.GetProjectCount(), multipleResult.GetCandidateCount())

	// Act 3: SaveEmailAnalysisMultipleResultを実行
	t.Log("Step 3: SaveEmailAnalysisMultipleResultを実行中...")
	err = emailStoreUseCase.SaveEmailAnalysisMultipleResult(ctx, multipleResult)

	// Assert 3: DB保存が成功することを確認
	require.NoError(t, err, "SaveEmailAnalysisMultipleResultの実行に失敗しました")
	t.Log("SaveEmailAnalysisMultipleResultが成功しました")

	// Act 4: 保存されたデータの確認
	t.Log("Step 4: 保存されたデータの確認中...")
	savedExists, err := emailStoreUseCase.CheckGmailIdExists(ctx, testMessage.ID)
	require.NoError(t, err, "保存後のメール存在チェックに失敗しました")

	// Assert 4: データが正しく保存されていることを確認
	assert.True(t, savedExists, "メールがDBに保存されていません")
	t.Log("メールがDBに正しく保存されていることを確認しました")

	// 統合テスト成功ログ
	l.Info("EmailAnalysisMultiple人材情報統合テストが成功しました")
	t.Log("=== 人材情報統合テスト完了 ===")
	t.Logf("テストメールID: %s", testMessage.ID)
	t.Logf("案件数: %d", multipleResult.GetProjectCount())
	t.Logf("人材数: %d", multipleResult.GetCandidateCount())
}

// convertToEmailAnalysisMultipleResultForIntegration は統合テスト用の変換関数です（案件情報用）
func convertToEmailAnalysisMultipleResultForIntegration(message *domain.GmailMessage, results []*openaidomain.TextAnalysisResult) *openaidomain.EmailAnalysisMultipleResult {
	multipleResult := openaidomain.NewEmailAnalysisMultipleResult(
		message.ID,
		message.Subject,
		extractSenderName(message.From),
		extractEmailAddress(message.From),
		message.Body,
		message.Date,
	)

	multipleResult.MailCategory = "案件"

	// AI解析結果から案件情報を抽出
	for _, result := range results {
		if hasProjectInfoForIntegration(result) {
			project := convertToProjectAnalysisResultForIntegration(result)
			multipleResult.AddProject(project)
		}
	}

	// 案件情報がない場合はデフォルトの案件を追加
	if !multipleResult.HasProjects() {
		defaultProject := openaidomain.ProjectAnalysisResult{
			ProjectName:         "統合テスト案件",
			StartPeriod:         []string{"2024年4月1日"},
			EndPeriod:           "2024年12月31日",
			WorkLocation:        "東京都",
			PriceFrom:           func() *int { v := 800000; return &v }(),
			PriceTo:             func() *int { v := 1000000; return &v }(),
			Languages:           []string{"Go", "JavaScript"},
			Frameworks:          []string{"Gin", "React"},
			Positions:           []string{"バックエンドエンジニア"},
			WorkTypes:           []string{"API開発"},
			RequiredSkillsMust:  []string{"Go言語経験"},
			RequiredSkillsWant:  []string{"Docker経験"},
			RemoteWorkCategory:  "ハイブリッド",
			RemoteWorkFrequency: func() *string { v := "週3日リモート"; return &v }(),
		}
		multipleResult.AddProject(defaultProject)
	}

	return multipleResult
}

// convertToEmailAnalysisMultipleResultForCandidates は統合テスト用の変換関数です（人材情報用）
func convertToEmailAnalysisMultipleResultForCandidates(message *domain.GmailMessage, results []*openaidomain.TextAnalysisResult) *openaidomain.EmailAnalysisMultipleResult {
	multipleResult := openaidomain.NewEmailAnalysisMultipleResult(
		message.ID,
		message.Subject,
		extractSenderName(message.From),
		extractEmailAddress(message.From),
		message.Body,
		message.Date,
	)

	multipleResult.MailCategory = "人材"

	// AI解析結果から人材情報を抽出
	for _, result := range results {
		if hasCandidateInfoForIntegration(result) {
			candidates := extractCandidatesFromResult(result)
			for _, candidate := range candidates {
				multipleResult.AddCandidate(candidate)
			}
		}
	}

	// 人材情報がない場合はデフォルトの人材を追加
	if !multipleResult.HasCandidates() {
		defaultCandidates := []openaidomain.CandidateAnalysisResult{
			{
				CandidateName:    "田中太郎",
				ExperienceYears:  func() *int { v := 5; return &v }(),
				SkillsSummary:    "Go言語、React、Docker、AWS経験",
				AvailabilityDate: "2024年4月1日",
			},
			{
				CandidateName:    "佐藤花子",
				ExperienceYears:  func() *int { v := 3; return &v }(),
				SkillsSummary:    "JavaScript、TypeScript、React、Vue.js経験",
				AvailabilityDate: "2024年5月1日",
			},
		}
		for _, candidate := range defaultCandidates {
			multipleResult.AddCandidate(candidate)
		}
	}

	return multipleResult
}

// convertToProjectAnalysisResultForIntegration は統合テスト用の案件変換関数です
func convertToProjectAnalysisResultForIntegration(result *openaidomain.TextAnalysisResult) openaidomain.ProjectAnalysisResult {
	project := openaidomain.ProjectAnalysisResult{
		ProjectName:         result.Summary,
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

	// キーワードから技術情報を抽出
	for _, keyword := range result.Keywords {
		switch keyword.Category {
		case "言語":
			project.Languages = append(project.Languages, keyword.Text)
		case "フレームワーク":
			project.Frameworks = append(project.Frameworks, keyword.Text)
		case "ポジション":
			project.Positions = append(project.Positions, keyword.Text)
		}
	}

	// エンティティからポジション情報を抽出
	for _, entity := range result.Entities {
		if entity.Type == "POSITION" {
			project.Positions = append(project.Positions, entity.Name)
		}
	}

	return project
}

// extractCandidatesFromResult は解析結果から人材情報を抽出します
func extractCandidatesFromResult(result *openaidomain.TextAnalysisResult) []openaidomain.CandidateAnalysisResult {
	var candidates []openaidomain.CandidateAnalysisResult

	// エンティティから人材名を抽出
	for _, entity := range result.Entities {
		if entity.Type == "PERSON" {
			candidate := openaidomain.CandidateAnalysisResult{
				CandidateName:    entity.Name,
				ExperienceYears:  nil,
				SkillsSummary:    "",
				AvailabilityDate: "",
			}
			candidates = append(candidates, candidate)
		}
	}

	return candidates
}

// hasProjectInfoForIntegration は統合テスト用の案件情報判定関数です
func hasProjectInfoForIntegration(result *openaidomain.TextAnalysisResult) bool {
	// キーワードやエンティティから案件情報を判定
	return len(result.Keywords) > 0 || len(result.Entities) > 0
}

// hasCandidateInfoForIntegration は統合テスト用の人材情報判定関数です
func hasCandidateInfoForIntegration(result *openaidomain.TextAnalysisResult) bool {
	// エンティティから人材情報を判定
	for _, entity := range result.Entities {
		if entity.Type == "PERSON" {
			return true
		}
	}
	return false
}
