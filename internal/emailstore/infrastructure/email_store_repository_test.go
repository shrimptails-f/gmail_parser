//go:build integration

// Package infrastructure はメール保存機能のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	"business/internal/emailstore/domain"
	openaidomain "business/internal/openai/domain"
	"business/tools/migrations/model"
	"business/tools/mysql"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailStoreRepositoryImpl_SaveEmail(t *testing.T) {
	// テスト用DBの準備
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル作成
	err = db.DB.AutoMigrate(
		model.KeywordGroup{},
		model.KeyWord{},
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
	ctx := context.Background()

	tests := []struct {
		name          string
		input         *openaidomain.EmailAnalysisResult
		expectedError string
		setupData     func()
	}{
		{
			name: "正常系_新規メール保存成功",
			input: &openaidomain.EmailAnalysisResult{
				ID:                  "test-email-id-1",
				Subject:             "テスト件名",
				From:                "sender@example.com",
				FromEmail:           "sender@example.com",
				Date:                time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Body:                "テスト本文",
				MailCategory:        "案件",
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
				RemoteWorkCategory:  "フルリモート",
				RemoteWorkFrequency: stringPtr("週5日"),
			},
			expectedError: "",
			setupData:     func() {},
		},
		{
			name: "正常系_案件以外のメール保存成功",
			input: &openaidomain.EmailAnalysisResult{
				ID:           "test-email-id-2",
				Subject:      "営業メール",
				From:         "sales@example.com",
				FromEmail:    "sales@example.com",
				Date:         time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
				Body:         "営業メール本文",
				MailCategory: "営業",
			},
			expectedError: "",
			setupData:     func() {},
		},
		{
			name: "異常系_重複メールID",
			input: &openaidomain.EmailAnalysisResult{
				ID:        "test-email-id-1",
				Subject:   "重複テスト",
				From:      "test@example.com",
				FromEmail: "test@example.com",
				Date:      time.Now(),
				Body:      "重複テスト本文",
			},
			expectedError: "メールが既に存在します",
			setupData: func() {
				// 事前に同じIDのメールを保存
				body := "既存メール本文"
				email := &domain.Email{
					ID:           "test-email-id-1",
					Subject:      "既存メール",
					SenderName:   "existing@example.com",
					SenderEmail:  "existing@example.com",
					ReceivedDate: time.Now(),
					Body:         &body,
					Category:     "案件",
				}
				db.DB.Create(email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupData()

			// Act
			err := repo.SaveEmail(ctx, tt.input)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)

				// データが正しく保存されているか確認
				var savedEmail domain.Email
				result := db.DB.Where("id = ?", tt.input.ID).First(&savedEmail)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.input.Subject, savedEmail.Subject)
				assert.Equal(t, tt.input.From, savedEmail.SenderName)
				assert.Equal(t, tt.input.FromEmail, savedEmail.SenderEmail)
				assert.Equal(t, tt.input.Body, *savedEmail.Body)

				// 案件メールの場合、EmailProjectも確認
				if tt.input.MailCategory == "案件" {
					var savedProject domain.EmailProject
					result := db.DB.Where("email_id = ?", tt.input.ID).First(&savedProject)
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
						var entryTimings []domain.EntryTiming
						result := db.DB.Where("email_project_id = ?", savedProject.EmailID).Find(&entryTimings)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.StartPeriod), len(entryTimings))
					}

					// キーワード関連の確認
					if len(tt.input.Languages) > 0 || len(tt.input.Frameworks) > 0 || len(tt.input.RequiredSkillsMust) > 0 || len(tt.input.RequiredSkillsWant) > 0 {
						var emailKeywordGroups []domain.EmailKeywordGroup
						result := db.DB.Where("email_id = ?", tt.input.ID).Find(&emailKeywordGroups)
						assert.NoError(t, result.Error)
						expectedKeywordCount := len(tt.input.Languages) + len(tt.input.Frameworks) + len(tt.input.RequiredSkillsMust) + len(tt.input.RequiredSkillsWant)
						assert.Equal(t, expectedKeywordCount, len(emailKeywordGroups))
					}

					// ポジション関連の確認
					if len(tt.input.Positions) > 0 {
						var emailPositionGroups []domain.EmailPositionGroup
						result := db.DB.Where("email_id = ?", tt.input.ID).Find(&emailPositionGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.Positions), len(emailPositionGroups))

						// PositionGroupとPositionWordが作成されているか確認
						for _, position := range tt.input.Positions {
							var positionGroup domain.PositionGroup
							result := db.DB.Where("name = ?", position).First(&positionGroup)
							assert.NoError(t, result.Error)

							var positionWord domain.PositionWord
							result = db.DB.Where("position_group_id = ? AND word = ?", positionGroup.PositionGroupID, position).First(&positionWord)
							assert.NoError(t, result.Error)
						}
					}

					// 業務種別関連の確認
					if len(tt.input.WorkTypes) > 0 {
						var emailWorkTypeGroups []domain.EmailWorkTypeGroup
						result := db.DB.Where("email_id = ?", tt.input.ID).Find(&emailWorkTypeGroups)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.WorkTypes), len(emailWorkTypeGroups))

						// WorkTypeGroupとWorkTypeWordが作成されているか確認
						for _, workType := range tt.input.WorkTypes {
							var workTypeGroup domain.WorkTypeGroup
							result := db.DB.Where("name = ?", workType).First(&workTypeGroup)
							assert.NoError(t, result.Error)

							var workTypeWord domain.WorkTypeWord
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

func TestEmailStoreRepositoryImpl_GetEmailByID(t *testing.T) {
	// テスト用DBの準備
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル作成
	err = db.DB.AutoMigrate(&domain.Email{})
	require.NoError(t, err)

	repo := NewEmailStoreRepository(db.DB)
	ctx := context.Background()

	// テストデータの準備
	body := "テスト本文"
	testEmail := &domain.Email{
		ID:           "test-email-id",
		Subject:      "テスト件名",
		SenderName:   "test@example.com",
		SenderEmail:  "test@example.com",
		ReceivedDate: time.Now(),
		Body:         &body,
		Category:     "案件",
	}
	db.DB.Create(testEmail)

	tests := []struct {
		name          string
		emailID       string
		expectedError string
	}{
		{
			name:          "正常系_メール取得成功",
			emailID:       "test-email-id",
			expectedError: "",
		},
		{
			name:          "異常系_存在しないメール",
			emailID:       "non-existent-id",
			expectedError: "メールが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			email, err := repo.GetEmailByID(ctx, tt.emailID)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, email)
				assert.Equal(t, tt.emailID, email.ID)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, email)
			}
		})
	}
}

func TestEmailStoreRepositoryImpl_EmailExists(t *testing.T) {
	// テスト用DBの準備
	db, cleanup, err := mysql.CreateNewTestDB()
	require.NoError(t, err)
	defer cleanup()

	// テーブル作成
	err = db.DB.AutoMigrate(&domain.Email{})
	require.NoError(t, err)

	repo := NewEmailStoreRepository(db.DB)
	ctx := context.Background()

	// テストデータの準備
	body := "既存メール本文"
	testEmail := &domain.Email{
		ID:           "existing-email-id",
		Subject:      "既存メール",
		SenderName:   "test@example.com",
		SenderEmail:  "test@example.com",
		ReceivedDate: time.Now(),
		Body:         &body,
		Category:     "案件",
	}
	db.DB.Create(testEmail)

	tests := []struct {
		name     string
		emailID  string
		expected bool
	}{
		{
			name:     "正常系_メール存在",
			emailID:  "existing-email-id",
			expected: true,
		},
		{
			name:     "正常系_メール存在しない",
			emailID:  "non-existent-id",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			exists, err := repo.EmailExists(ctx, tt.emailID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

// intPtr はintのポインタを返すヘルパー関数です
func intPtr(i int) *int {
	return &i
}

// stringPtr はstringのポインタを返すヘルパー関数です
func stringPtr(s string) *string {
	return &s
}
