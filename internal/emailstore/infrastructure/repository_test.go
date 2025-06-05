// // Package infrastructure はメール保存機能のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	"business/tools/migrations/model"
	"business/tools/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailStoreRepositoryImpl_SaveEmail(t *testing.T) {
	t.Parallel()

	// テスト用DB初期化
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル準備
	err = db.DB.AutoMigrate(
		model.KeywordGroup{},
		model.KeyWord{},
		model.KeywordGroupWordLink{},
		model.PositionGroup{},
		model.PositionWord{},
		model.WorkTypeGroup{},
		model.WorkTypeWord{},
		model.Email{},
		model.EmailProject{},
		model.EmailCandidate{},
		model.EntryTiming{},
		model.EmailKeywordGroup{},
		model.EmailPositionGroup{},
		model.EmailWorkTypeGroup{},
	)
	require.NoError(t, err)

	repo := &EmailStoreRepositoryImpl{db: db.DB}

	now := time.Now()

	email := Email{
		GmailID:      "dummy-gmail-id-12345",
		Subject:      "テスト件名",
		SenderName:   "田中 太郎",
		SenderEmail:  "tanaka@example.com",
		ReceivedDate: time.Date(2024, 4, 15, 10, 0, 0, 0, time.UTC),
		Body:         stringPtr("これはテスト本文です"),
		Category:     "案件情報",
		EmailProject: &EmailProject{
			ProjectTitle:    stringPtr("プロジェクトZ"),
			EntryTiming:     stringPtr("2024年7月"),
			EndTiming:       stringPtr("2025年1月"),
			WorkLocation:    stringPtr("東京都 千代田区"),
			PriceFrom:       intPtr(600000),
			PriceTo:         intPtr(800000),
			Languages:       stringPtr("Go, Python"),
			Frameworks:      stringPtr("Gin, React"),
			Positions:       stringPtr("バックエンドエンジニア, インフラエンジニア"),
			WorkTypes:       stringPtr("Web開発, 保守運用"),
			MustSkills:      stringPtr("Docker, AWS"),
			WantSkills:      stringPtr("Kubernetes, Terraform"),
			RemoteType:      stringPtr("フルリモート"),
			RemoteFrequency: stringPtr("週5"),
		},
		EntryTimings: []EntryTiming{
			{
				StartDate: "６月上旬",
			},
		},
		EmailKeywordGroups: []EmailKeywordGroup{
			{
				KeywordGroup: KeywordGroup{
					Name: "Go",
					Type: "language",
					WordLinks: []KeywordGroupWordLink{
						{
							KeyWord: KeyWord{
								Word: "Go",
							},
						},
					},
				},
			},
		},
		EmailPositionGroups: []EmailPositionGroup{
			{
				PositionGroup: PositionGroup{
					Name:      "バックエンドエンジニア",
					CreatedAt: now,
					UpdatedAt: now,
					Words: []PositionWord{
						{
							Word:      "BE開発",
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
				},
			},
		},
		EmailWorkTypeGroups: []EmailWorkTypeGroup{
			{
				WorkTypeGroup: WorkTypeGroup{
					Name:      "開発",
					CreatedAt: now,
					UpdatedAt: now,
					Words: []WorkTypeWord{
						{
							Word:      "Web開発",
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
				},
			},
		},
	}

	err = repo.SaveEmail(email)
	assert.NoError(t, err)

	// 登録確認（emails）
	var actualEmail Email
	err = db.DB.
		Preload("EmailProject").
		Preload("EntryTimings").
		Preload("EmailKeywordGroups.KeywordGroup.WordLinks.KeyWord").
		Preload("EmailPositionGroups.PositionGroup.Words").
		Preload("EmailWorkTypeGroups.WorkTypeGroup.Words").
		First(&actualEmail, "gmail_id = ?", "dummy-gmail-id-12345").
		Error

	assert.NoError(t, err)
	assert.Equal(t, "テスト件名", actualEmail.Subject)
	assert.Equal(t, "プロジェクトZ", *actualEmail.EmailProject.ProjectTitle)
	assert.Equal(t, "６月上旬", actualEmail.EntryTimings[0].StartDate)
	assert.Equal(t, "Go", actualEmail.EmailKeywordGroups[0].KeywordGroup.Name)
	assert.Equal(t, "バックエンドエンジニア", actualEmail.EmailPositionGroups[0].PositionGroup.Name)
	assert.Equal(t, "開発", actualEmail.EmailWorkTypeGroups[0].WorkTypeGroup.Name)

}
func stringPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }
