// Package di はメール保存機能の依存性注入の統合テストを提供します。
package di

// import (
// 	"business/internal/emailstore/domain"
// 	openaidomain "business/internal/openai/domain"
// 	"business/tools/migrations/model"
// 	"business/tools/mysql"
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // func TestEmailStoreIntegration_SaveEmailAnalysisResult(t *testing.T) {
// // 	// テスト用DBの準備
// // 	db, cleanup, err := mysql.CreateNewTestDB()
// // 	require.NoError(t, err)
// // 	defer cleanup()

// // 	// マイグレーションファイルと同じ順序でテーブル作成
// // 	err = db.DB.AutoMigrate(
// // 		&model.KeywordGroup{},
// // 		&model.KeyWord{},
// // 		&model.KeywordGroupWordLink{},
// // 		&model.PositionGroup{},
// // 		&model.PositionWord{},
// // 		&model.WorkTypeGroup{},
// // 		&model.WorkTypeWord{},
// // 		&model.Email{},
// // 		&model.EmailProject{},
// // 		&model.EmailCandidate{},
// // 		&model.EntryTiming{},
// // 		&model.EmailKeywordGroup{},
// // 		&model.EmailPositionGroup{},
// // 		&model.EmailWorkTypeGroup{},
// // 	)
// // 	require.NoError(t, err)

// // 	// DIコンテナからユースケースを取得
// // 	emailStoreUseCase := ProvideEmailStoreDependencies(db.DB)
// // 	ctx := context.Background()

// // 	tests := []struct {
// // 		name          string
// // 		input         *openaidomain.EmailAnalysisResult
// // 		expectedError string
// // 		verify        func(t *testing.T)
// // 	}{
// // 		{
// // 			name: "統合テスト_案件メール保存成功",
// // 			input: &openaidomain.EmailAnalysisResult{
// // 				GmailID:             "integration-test-email-1",
// // 				Subject:             "【案件】Go言語エンジニア募集",
// // 				From:                "recruiter@example.com",
// // 				FromEmail:           "recruiter@example.com",
// // 				Date:                time.Date(2024, 5, 30, 12, 0, 0, 0, time.UTC),
// // 				Body:                "Go言語を使用したバックエンド開発の案件です。",
// // 				MailCategory:        "案件",
// // 				StartPeriod:         []string{"2024年6月", "2024年7月"},
// // 				EndPeriod:           "2024年12月",
// // 				WorkLocation:        "東京都渋谷区",
// // 				PriceFrom:           intPtr(600000),
// // 				PriceTo:             intPtr(800000),
// // 				Languages:           []string{"Go", "TypeScript"},
// // 				Frameworks:          []string{"Gin", "React"},
// // 				RequiredSkillsMust:  []string{"Git", "Docker", "AWS"},
// // 				RequiredSkillsWant:  []string{"Kubernetes", "Terraform"},
// // 				RemoteWorkCategory:  "ハイブリッド",
// // 				RemoteWorkFrequency: stringPtr("週3日"),
// // 			},
// // 			expectedError: "",
// // 			verify: func(t *testing.T) {
// // 				// Emailテーブルの確認
// // 				var savedEmail domain.Email
// // 				result := db.DB.Where("gmail_id = ?", "integration-test-email-1").First(&savedEmail)
// // 				assert.NoError(t, result.Error)
// // 				assert.Equal(t, "【案件】Go言語エンジニア募集", savedEmail.Subject)

// // 				// EmailProjectテーブルの確認
// // 				var savedProject domain.EmailProject
// // 				result = db.DB.Where("email_id = ?", savedEmail.ID).First(&savedProject)
// // 				assert.NoError(t, result.Error)
// // 				assert.Equal(t, "東京都渋谷区", *savedProject.WorkLocation)
// // 				assert.Equal(t, 600000, *savedProject.PriceFrom)
// // 				assert.Equal(t, 800000, *savedProject.PriceTo)
// // 				assert.Equal(t, "ハイブリッド", *savedProject.RemoteType)
// // 				assert.Equal(t, "週3日", *savedProject.RemoteFrequency)

// // 				// EntryTimingテーブルの確認
// // 				var entryTimings []domain.EntryTiming
// // 				result = db.DB.Where("email_project_id = ?", savedProject.EmailID).Find(&entryTimings)
// // 				assert.NoError(t, result.Error)
// // 				assert.Len(t, entryTimings, 2)
// // 				timings := make([]string, len(entryTimings))
// // 				for i, timing := range entryTimings {
// // 					timings[i] = timing.StartDate
// // 				}
// // 				assert.Contains(t, timings, "2024年6月")
// // 				assert.Contains(t, timings, "2024年7月")

// // 				// KeywordGroupとEmailKeywordGroupの確認
// // 				var keywordGroups []domain.KeywordGroup
// // 				result = db.DB.Find(&keywordGroups)
// // 				assert.NoError(t, result.Error)
// // 				assert.GreaterOrEqual(t, len(keywordGroups), 9) // 最低9個のキーワードグループが作成される

// // 				// EmailKeywordGroupの確認
// // 				var emailKeywordGroups []domain.EmailKeywordGroup
// // 				result = db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailKeywordGroups)
// // 				assert.NoError(t, result.Error)
// // 				assert.Len(t, emailKeywordGroups, 9) // 9個のキーワード関連付け

// // 				// タイプ別の確認
// // 				typeCount := make(map[string]int)
// // 				for _, ekg := range emailKeywordGroups {
// // 					typeCount[ekg.Type]++
// // 				}
// // 				assert.Equal(t, 2, typeCount["LANGUAGE"])  // Go, TypeScript
// // 				assert.Equal(t, 2, typeCount["FRAMEWORK"]) // Gin, React
// // 				assert.Equal(t, 3, typeCount["MUST"])      // Git, Docker, AWS
// // 				assert.Equal(t, 2, typeCount["WANT"])      // Kubernetes, Terraform
// // 			},
// // 		},
// // 		{
// // 			name: "統合テスト_営業メール保存成功",
// // 			input: &openaidomain.EmailAnalysisResult{
// // 				GmailID:      "integration-test-email-2",
// // 				Subject:      "営業のご案内",
// // 				From:         "sales@example.com",
// // 				FromEmail:    "sales@example.com",
// // 				Date:         time.Date(2024, 5, 30, 13, 0, 0, 0, time.UTC),
// // 				Body:         "弊社サービスのご案内です。",
// // 				MailCategory: "営業",
// // 			},
// // 			expectedError: "",
// // 			verify: func(t *testing.T) {
// // 				// Emailテーブルの確認
// // 				var savedEmail domain.Email
// // 				result := db.DB.Where("gmail_id = ?", "integration-test-email-2").First(&savedEmail)
// // 				assert.NoError(t, result.Error)
// // 				assert.Equal(t, "営業のご案内", savedEmail.Subject)

// // 				// EmailProjectテーブルには保存されないことを確認
// // 				var projectCount int64
// // 				result = db.DB.Model(&domain.EmailProject{}).Where("email_id = ?", savedEmail.ID).Count(&projectCount)
// // 				assert.NoError(t, result.Error)
// // 				assert.Equal(t, int64(0), projectCount)
// // 			},
// // 		},
// // 	}

// // 	for _, tt := range tests {
// // 		t.Run(tt.name, func(t *testing.T) {
// // 			// Act
// // 			err := emailStoreUseCase.SaveEmailAnalysisResult(ctx, tt.input)

// // 			// Assert
// // 			if tt.expectedError == "" {
// // 				assert.NoError(t, err)
// // 				if tt.verify != nil {
// // 					tt.verify(t)
// // 				}
// // 			} else {
// // 				assert.Error(t, err)
// // 				assert.Contains(t, err.Error(), tt.expectedError)
// // 			}
// // 		})
// // 	}
// // }

// // // intPtr はintのポインタを返すヘルパー関数です
// // func intPtr(i int) *int {
// // 	return &i
// // }

// // // stringPtr はstringのポインタを返すヘルパー関数です
// // func stringPtr(s string) *string {
// // 	return &s
// // }
