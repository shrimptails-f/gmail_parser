// // Package infrastructure はメール保存機能のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	cd "business/internal/common/domain"
	"business/tools/migrations/model"
	"business/tools/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailStoreRepositoryImpl_SaveEmail(t *testing.T) {
	t.Parallel()
	// テスト用DBの準備
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル作成
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

	repo := NewEmailStoreRepository(db.DB)

	tests := []struct {
		name          string
		input         cd.Email
		expectedError string
		setupData     func()
	}{
		{
			name: "正常系_新規メール保存成功",
			input: cd.Email{
				GmailID:             "test-email-id-1",
				Subject:             "テスト件名",
				From:                "sender@example.com",
				FromEmail:           "sender@example.com",
				ReceivedDate:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Body:                "テスト本文",
				Category:            "案件",
				StartPeriod:         []string{"2024年4月", "2024年5月"},
				EndPeriod:           "2024年12月",
				WorkLocation:        "東京都",
				PriceFrom:           intPtr(500000),
				PriceTo:             intPtr(600000),
				Languages:           []string{"Go", "Python"},
				Frameworks:          []string{"Gin", "Django"},
				Positions:           []string{"PM", "SE"},
				WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
				RequiredSkillsMust:  []string{"Git", "Docker"},
				RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
				RemoteWorkCategory:  stringPtr("フルリモート"),
				RemoteWorkFrequency: stringPtr("週5日"),
			},
			expectedError: "",
			setupData:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupData()

			// Act
			err := repo.SaveEmail(tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)

				// データが正しく保存されているか確認
				var savedEmail Email
				result := db.DB.Where("gmail_id = ?", tt.input.GmailID).First(&savedEmail)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.input.Subject, savedEmail.Subject)
				assert.Equal(t, tt.input.From, savedEmail.SenderName)
				assert.Equal(t, tt.input.FromEmail, savedEmail.SenderEmail)
				assert.Equal(t, tt.input.Body, *savedEmail.Body)

				// 案件メールの場合、EmailProjectも確認
				if tt.input.Category == "案件" {
					var savedProject EmailProject
					result := db.DB.Where("email_id = ?", savedEmail.ID).First(&savedProject)
					assert.NoError(t, result.Error)
					if savedProject.EndTiming != nil {
						assert.Equal(t, tt.input.EndPeriod, *savedProject.EndTiming)
					}
					if savedProject.WorkLocation != nil {
						assert.Equal(t, tt.input.WorkLocation, *savedProject.WorkLocation)
					}

					if tt.input.PriceFrom != nil {
						assert.Equal(t, *tt.input.PriceFrom, *savedProject.PriceFrom)
					}
					if tt.input.PriceTo != nil {
						assert.Equal(t, *tt.input.PriceTo, *savedProject.PriceTo)
					}

					// EntryTimingの確認
					if len(tt.input.StartPeriod) > 0 {
						var entryTimings []EntryTiming
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&entryTimings)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.StartPeriod), len(entryTimings))
					}

					// キーワード関連の確認
					if len(tt.input.Languages) > 0 || len(tt.input.Frameworks) > 0 || len(tt.input.RequiredSkillsMust) > 0 || len(tt.input.RequiredSkillsWant) > 0 {
						var emailKeywordGroups []EmailKeywordGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailKeywordGroups)
						assert.NoError(t, result.Error)
						expectedKeywordCount := len(tt.input.Languages) + len(tt.input.Frameworks) + len(tt.input.RequiredSkillsMust) + len(tt.input.RequiredSkillsWant)
						assert.Equal(t, expectedKeywordCount, len(emailKeywordGroups))
					}

					// ポジション関連の確認
					if len(tt.input.Positions) > 0 {
						var emailPositionGroups []EmailPositionGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailPositionGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.Positions), len(emailPositionGroups))

						// PositionGroupとPositionWordが作成されているか確認
						for _, position := range tt.input.Positions {
							var positionGroup PositionGroup
							result := db.DB.Where("name = ?", position).First(&positionGroup)
							assert.NoError(t, result.Error)

							var positionWord PositionWord
							result = db.DB.Where("position_group_id = ? AND word = ?", positionGroup.PositionGroupID, position).First(&positionWord)
							assert.NoError(t, result.Error)
						}
					}

					// 業務種別関連の確認
					if len(tt.input.WorkTypes) > 0 {
						var emailWorkTypeGroups []EmailWorkTypeGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailWorkTypeGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.WorkTypes), len(emailWorkTypeGroups))

						// WorkTypeGroupとWorkTypeWordが作成されているか確認
						for _, workType := range tt.input.WorkTypes {
							var workTypeGroup WorkTypeGroup
							result := db.DB.Where("name = ?", workType).First(&workTypeGroup)
							assert.NoError(t, result.Error)

							var workTypeWord WorkTypeWord
							result = db.DB.Where("work_type_group_id = ? AND word = ?", workTypeGroup.WorkTypeGroupID, workType).First(&workTypeWord)
							assert.NoError(t, result.Error)
						}
					}
				}
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
func TestEmailStoreRepositoryImpl_SavedDuplicateEmail(t *testing.T) {
	t.Parallel()
	// テスト用DBの準備
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル作成
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

	repo := NewEmailStoreRepository(db.DB)

	tests := []struct {
		name          string
		input         cd.Email
		expectedError string
		setupData     func()
	}{
		{
			name: "同じような内容のメールが保存できること 1通目",
			input: cd.Email{
				GmailID:             "test-email-id-1",
				Subject:             "テスト件名",
				From:                "sender@example.com",
				FromEmail:           "sender@example.com",
				ReceivedDate:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Body:                "テスト本文",
				Category:            "案件",
				StartPeriod:         []string{"2024年4月", "2024年5月"},
				EndPeriod:           "2024年12月",
				WorkLocation:        "東京都",
				PriceFrom:           intPtr(500000),
				PriceTo:             intPtr(600000),
				Languages:           []string{"Go", "Python"},
				Frameworks:          []string{"Gin", "Django"},
				Positions:           []string{"PM", "SE"},
				WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
				RequiredSkillsMust:  []string{"Git", "Docker"},
				RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
				RemoteWorkCategory:  stringPtr("フルリモート"),
				RemoteWorkFrequency: stringPtr("週5日"),
			},
			expectedError: "",
			setupData:     func() {},
		},
		{
			name: "同じような内容のメールが保存できること ２通目",
			input: cd.Email{
				GmailID:             "test-email-id-２",
				Subject:             "テスト件名",
				From:                "sender@example.com",
				FromEmail:           "sender@example.com",
				ReceivedDate:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Body:                "テスト本文",
				Category:            "案件",
				StartPeriod:         []string{"2024年4月", "2024年5月"},
				EndPeriod:           "2024年12月",
				WorkLocation:        "東京都",
				PriceFrom:           intPtr(500000),
				PriceTo:             intPtr(600000),
				Languages:           []string{"Go", "Python"},
				Frameworks:          []string{"Gin", "Django"},
				Positions:           []string{"PM", "SE"},
				WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
				RequiredSkillsMust:  []string{"Git", "Docker"},
				RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
				RemoteWorkCategory:  stringPtr("フルリモート"),
				RemoteWorkFrequency: stringPtr("週5日"),
			},
			expectedError: "",
			setupData:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupData()

			// Act
			err := repo.SaveEmail(tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)

				// データが正しく保存されているか確認
				var savedEmail Email
				result := db.DB.Where("gmail_id = ?", tt.input.GmailID).First(&savedEmail)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.input.Subject, savedEmail.Subject)
				assert.Equal(t, tt.input.From, savedEmail.SenderName)
				assert.Equal(t, tt.input.FromEmail, savedEmail.SenderEmail)
				assert.Equal(t, tt.input.Body, *savedEmail.Body)

				// 案件メールの場合、EmailProjectも確認
				if tt.input.Category == "案件" {
					var savedProject EmailProject
					result := db.DB.Where("email_id = ?", savedEmail.ID).First(&savedProject)
					assert.NoError(t, result.Error)
					if savedProject.EndTiming != nil {
						assert.Equal(t, tt.input.EndPeriod, *savedProject.EndTiming)
					}
					if savedProject.WorkLocation != nil {
						assert.Equal(t, tt.input.WorkLocation, *savedProject.WorkLocation)
					}

					if tt.input.PriceFrom != nil {
						assert.Equal(t, *tt.input.PriceFrom, *savedProject.PriceFrom)
					}
					if tt.input.PriceTo != nil {
						assert.Equal(t, *tt.input.PriceTo, *savedProject.PriceTo)
					}

					// EntryTimingの確認
					if len(tt.input.StartPeriod) > 0 {
						var entryTimings []EntryTiming
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&entryTimings)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.StartPeriod), len(entryTimings))
					}

					// キーワード関連の確認
					if len(tt.input.Languages) > 0 || len(tt.input.Frameworks) > 0 || len(tt.input.RequiredSkillsMust) > 0 || len(tt.input.RequiredSkillsWant) > 0 {
						var emailKeywordGroups []EmailKeywordGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailKeywordGroups)
						assert.NoError(t, result.Error)
						expectedKeywordCount := len(tt.input.Languages) + len(tt.input.Frameworks) + len(tt.input.RequiredSkillsMust) + len(tt.input.RequiredSkillsWant)
						assert.Equal(t, expectedKeywordCount, len(emailKeywordGroups))
					}

					// ポジション関連の確認
					if len(tt.input.Positions) > 0 {
						var emailPositionGroups []EmailPositionGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailPositionGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.Positions), len(emailPositionGroups))

						// PositionGroupとPositionWordが作成されているか確認
						for _, position := range tt.input.Positions {
							var positionGroup PositionGroup
							result := db.DB.Where("name = ?", position).First(&positionGroup)
							assert.NoError(t, result.Error)

							var positionWord PositionWord
							result = db.DB.Where("position_group_id = ? AND word = ?", positionGroup.PositionGroupID, position).First(&positionWord)
							assert.NoError(t, result.Error)
						}
					}

					// 業務種別関連の確認
					if len(tt.input.WorkTypes) > 0 {
						var emailWorkTypeGroups []EmailWorkTypeGroup
						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailWorkTypeGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.WorkTypes), len(emailWorkTypeGroups))

						// WorkTypeGroupとWorkTypeWordが作成されているか確認
						for _, workType := range tt.input.WorkTypes {
							var workTypeGroup WorkTypeGroup
							result := db.DB.Where("name = ?", workType).First(&workTypeGroup)
							assert.NoError(t, result.Error)

							var workTypeWord WorkTypeWord
							result = db.DB.Where("work_type_group_id = ? AND word = ?", workTypeGroup.WorkTypeGroupID, workType).First(&workTypeWord)
							assert.NoError(t, result.Error)
						}
					}
				}
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func stringPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }
