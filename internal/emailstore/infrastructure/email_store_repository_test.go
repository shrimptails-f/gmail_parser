//go:build integration

// Package infrastructure はメール保存機能のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	"business/internal/emailstore/domain"
	openaidomain "business/internal/openai/domain"
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
		&domain.Email{},
		&domain.EmailProject{},
		&domain.EntryTiming{},
		&domain.KeywordGroup{},
		&domain.KeyWord{},
		&domain.EmailKeywordGroup{},
		&domain.PositionGroup{},
		&domain.PositionWord{},
		&domain.EmailPositionGroup{},
		&domain.WorkTypeGroup{},
		&domain.WorkTypeWord{},
		&domain.EmailWorkTypeGroup{},
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
				email := &domain.Email{
					ID:        "test-email-id-1",
					Subject:   "既存メール",
					From:      "existing@example.com",
					FromEmail: "existing@example.com",
					Date:      time.Now(),
					Body:      "既存メール本文",
				}
				db.DB.Create(email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストデータのクリーンアップ
			db.DB.Exec("DELETE FROM email_work_type_groups")
			db.DB.Exec("DELETE FROM email_position_groups")
			db.DB.Exec("DELETE FROM email_keyword_groups")
			db.DB.Exec("DELETE FROM entry_timings")
			db.DB.Exec("DELETE FROM email_projects")
			db.DB.Exec("DELETE FROM emails")
			db.DB.Exec("DELETE FROM work_type_words")
			db.DB.Exec("DELETE FROM work_type_groups")
			db.DB.Exec("DELETE FROM position_words")
			db.DB.Exec("DELETE FROM position_groups")
			db.DB.Exec("DELETE FROM key_words")
			db.DB.Exec("DELETE FROM keyword_groups")

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
				assert.Equal(t, tt.input.From, savedEmail.From)
				assert.Equal(t, tt.input.FromEmail, savedEmail.FromEmail)
				assert.Equal(t, tt.input.Body, savedEmail.Body)

				// 案件メールの場合、EmailProjectも確認
				if tt.input.MailCategory == "案件" {
					var savedProject domain.EmailProject
					result := db.DB.Where("email_id = ?", tt.input.ID).First(&savedProject)
					assert.NoError(t, result.Error)
					assert.Equal(t, tt.input.MailCategory, savedProject.MailCategory)
					assert.Equal(t, tt.input.EndPeriod, savedProject.EndPeriod)
					assert.Equal(t, tt.input.WorkLocation, savedProject.WorkLocation)

					if tt.input.PriceFrom != nil {
						assert.Equal(t, *tt.input.PriceFrom, *savedProject.PriceFrom)
					}
					if tt.input.PriceTo != nil {
						assert.Equal(t, *tt.input.PriceTo, *savedProject.PriceTo)
					}

					// EntryTimingの確認
					if len(tt.input.StartPeriod) > 0 {
						var entryTimings []domain.EntryTiming
						result := db.DB.Where("email_project_id = ?", savedProject.ID).Find(&entryTimings)
						assert.NoError(t, result.Error)
						assert.Equal(t, len(tt.input.StartPeriod), len(entryTimings))
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
	testEmail := &domain.Email{
		ID:        "test-email-id",
		Subject:   "テスト件名",
		From:      "test@example.com",
		FromEmail: "test@example.com",
		Date:      time.Now(),
		Body:      "テスト本文",
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
	testEmail := &domain.Email{
		ID:        "existing-email-id",
		Subject:   "既存メール",
		From:      "test@example.com",
		FromEmail: "test@example.com",
		Date:      time.Now(),
		Body:      "既存メール本文",
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
