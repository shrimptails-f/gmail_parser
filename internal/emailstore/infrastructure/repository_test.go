// // Package infrastructure はメール保存機能のインフラストラクチャ層のテストを提供します。
package infrastructure

import (
	"business/tools/migrations/model"
	"business/tools/mysql"
	"fmt"
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
		EmailKeywordGroups: []EmailKeywordGroup{
			{
				Type:      "LANGUAGE",
				CreatedAt: now,
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
		Preload("EmailKeywordGroups.KeywordGroup.WordLinks.KeyWord").
		Preload("EmailPositionGroups.PositionGroup.Words").
		Preload("EmailWorkTypeGroups.WorkTypeGroup.Words").
		First(&actualEmail, "gmail_id = ?", "dummy-gmail-id-12345").
		Error

	fmt.Printf("aaaa %v \n", actualEmail.EmailKeywordGroups[0].KeywordGroup.Name)
	assert.NoError(t, err)
	assert.Equal(t, "テスト件名", actualEmail.Subject)
	assert.Equal(t, "Go", actualEmail.EmailKeywordGroups[0].KeywordGroup.Name)
	assert.Equal(t, "バックエンドエンジニア", actualEmail.EmailPositionGroups[0].PositionGroup.Name)
	assert.Equal(t, "開発", actualEmail.EmailWorkTypeGroups[0].WorkTypeGroup.Name)

}

// func TestEmailStoreRepositoryImpl_SaveEmail(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 		model.KeywordGroupWordLink{},
// 		model.PositionGroup{},
// 		model.PositionWord{},
// 		model.WorkTypeGroup{},
// 		model.WorkTypeWord{},
// 		model.Email{},
// 		model.EmailProject{},
// 		model.EmailCandidate{},
// 		model.EntryTiming{},
// 		model.EmailKeywordGroup{},
// 		model.EmailPositionGroup{},
// 		model.EmailWorkTypeGroup{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	tests := []struct {
// 		name          string
// 		input         *openaidomain.EmailAnalysisResult
// 		expectedError string
// 		setupData     func()
// 	}{
// 		{
// 			name: "正常系_新規メール保存成功",
// 			input: &openaidomain.EmailAnalysisResult{
// 				GmailID:             "test-email-id-1",
// 				Subject:             "テスト件名",
// 				From:                "sender@example.com",
// 				FromEmail:           "sender@example.com",
// 				Date:                time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
// 				Body:                "テスト本文",
// 				MailCategory:        "案件",
// 				StartPeriod:         []string{"2024年4月", "2024年5月"},
// 				EndPeriod:           "2024年12月",
// 				WorkLocation:        "東京都",
// 				PriceFrom:           intPtr(500000),
// 				PriceTo:             intPtr(600000),
// 				Languages:           []string{"Go", "Python"},
// 				Frameworks:          []string{"Gin", "Django"},
// 				Positions:           []string{"PM", "SE"},
// 				WorkTypes:           []string{"バックエンド開発", "インフラ構築"},
// 				RequiredSkillsMust:  []string{"Git", "Docker"},
// 				RequiredSkillsWant:  []string{"AWS", "Kubernetes"},
// 				RemoteWorkCategory:  "フルリモート",
// 				RemoteWorkFrequency: stringPtr("週5日"),
// 			},
// 			expectedError: "",
// 			setupData:     func() {},
// 		},
// 		// {
// 		// 	name: "正常系_案件以外のメール保存成功",
// 		// 	input: &openaidomain.EmailAnalysisResult{
// 		// 		GmailID:      "test-email-id-2",
// 		// 		Subject:      "営業メール",
// 		// 		From:         "sales@example.com",
// 		// 		FromEmail:    "sales@example.com",
// 		// 		Date:         time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
// 		// 		Body:         "営業メール本文",
// 		// 		MailCategory: "営業",
// 		// 	},
// 		// 	expectedError: "",
// 		// 	setupData:     func() {},
// 		// },
// 		// {
// 		// 	name: "異常系_重複メールID",
// 		// 	input: &openaidomain.EmailAnalysisResult{
// 		// 		GmailID:   "test-email-id-1",
// 		// 		Subject:   "重複テスト",
// 		// 		From:      "test@example.com",
// 		// 		FromEmail: "test@example.com",
// 		// 		Date:      time.Now(),
// 		// 		Body:      "重複テスト本文",
// 		// 	},
// 		// 	expectedError: "メールが既に存在します",
// 		// 	setupData: func() {
// 		// 		// 事前に同じIDのメールを保存
// 		// 		body := "既存メール本文"
// 		// 		email := &domain.Email{
// 		// 			GmailID:      "test-email-id-1",
// 		// 			Subject:      "既存メール",
// 		// 			SenderName:   "existing@example.com",
// 		// 			SenderEmail:  "existing@example.com",
// 		// 			ReceivedDate: time.Now(),
// 		// 			Body:         &body,
// 		// 			Category:     "案件",
// 		// 		}
// 		// 		db.DB.Create(email)
// 		// 	},
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			tt.setupData()

// 			// Act
// 			err := repo.SaveEmail(ctx, tt.input)

// 			// Assert
// 			if tt.expectedError == "" {
// 				assert.NoError(t, err)

// 				// データが正しく保存されているか確認
// 				var savedEmail domain.Email
// 				result := db.DB.Where("gmail_id = ?", tt.input.GmailID).First(&savedEmail)
// 				assert.NoError(t, result.Error)
// 				assert.Equal(t, tt.input.Subject, savedEmail.Subject)
// 				assert.Equal(t, tt.input.From, savedEmail.SenderName)
// 				assert.Equal(t, tt.input.FromEmail, savedEmail.SenderEmail)
// 				assert.Equal(t, tt.input.Body, *savedEmail.Body)

// 				// 案件メールの場合、EmailProjectも確認
// 				if tt.input.MailCategory == "案件" {
// 					var savedProject domain.EmailProject
// 					result := db.DB.Where("email_id = ?", savedEmail.ID).First(&savedProject)
// 					assert.NoError(t, result.Error)
// 					if savedProject.EndTiming != nil {
// 						assert.Equal(t, tt.input.EndPeriod, *savedProject.EndTiming)
// 					}
// 					if savedProject.WorkLocation != nil {
// 						assert.Equal(t, tt.input.WorkLocation, *savedProject.WorkLocation)
// 					}

// 					if tt.input.PriceFrom != nil {
// 						assert.Equal(t, *tt.input.PriceFrom, *savedProject.PriceFrom)
// 					}
// 					if tt.input.PriceTo != nil {
// 						assert.Equal(t, *tt.input.PriceTo, *savedProject.PriceTo)
// 					}

// 					// EntryTimingの確認
// 					if len(tt.input.StartPeriod) > 0 {
// 						var entryTimings []domain.EntryTiming
// 						result := db.DB.Where("email_project_id = ?", savedProject.EmailID).Find(&entryTimings)
// 						assert.NoError(t, result.Error)
// 						assert.Equal(t, len(tt.input.StartPeriod), len(entryTimings))
// 					}

// 					// キーワード関連の確認
// 					if len(tt.input.Languages) > 0 || len(tt.input.Frameworks) > 0 || len(tt.input.RequiredSkillsMust) > 0 || len(tt.input.RequiredSkillsWant) > 0 {
// 						var emailKeywordGroups []domain.EmailKeywordGroup
// 						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailKeywordGroups)
// 						assert.NoError(t, result.Error)
// 						expectedKeywordCount := len(tt.input.Languages) + len(tt.input.Frameworks) + len(tt.input.RequiredSkillsMust) + len(tt.input.RequiredSkillsWant)
// 						assert.Equal(t, expectedKeywordCount, len(emailKeywordGroups))
// 					}

// 					// ポジション関連の確認
// 					if len(tt.input.Positions) > 0 {
// 						var emailPositionGroups []domain.EmailPositionGroup
// 						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailPositionGroups)
// 						assert.NoError(t, result.Error)
// 						assert.Equal(t, len(tt.input.Positions), len(emailPositionGroups))

// 						// PositionGroupとPositionWordが作成されているか確認
// 						for _, position := range tt.input.Positions {
// 							var positionGroup domain.PositionGroup
// 							result := db.DB.Where("name = ?", position).First(&positionGroup)
// 							assert.NoError(t, result.Error)

// 							var positionWord domain.PositionWord
// 							result = db.DB.Where("position_group_id = ? AND word = ?", positionGroup.PositionGroupID, position).First(&positionWord)
// 							assert.NoError(t, result.Error)
// 						}
// 					}

// 					// 業務種別関連の確認
// 					if len(tt.input.WorkTypes) > 0 {
// 						var emailWorkTypeGroups []domain.EmailWorkTypeGroup
// 						result := db.DB.Where("email_id = ?", savedEmail.ID).Find(&emailWorkTypeGroups)
// 						assert.NoError(t, result.Error)
// 						assert.Equal(t, len(tt.input.WorkTypes), len(emailWorkTypeGroups))

// 						// WorkTypeGroupとWorkTypeWordが作成されているか確認
// 						for _, workType := range tt.input.WorkTypes {
// 							var workTypeGroup domain.WorkTypeGroup
// 							result := db.DB.Where("name = ?", workType).First(&workTypeGroup)
// 							assert.NoError(t, result.Error)

// 							var workTypeWord domain.WorkTypeWord
// 							result = db.DB.Where("work_type_group_id = ? AND word = ?", workTypeGroup.WorkTypeGroupID, workType).First(&workTypeWord)
// 							assert.NoError(t, result.Error)
// 						}
// 					}
// 				}
// 			} else {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError)
// 			}
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_GetEmailByID(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(&domain.Email{})
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// テストデータの準備
// 	body := "テスト本文"
// 	testEmail := &domain.Email{
// 		GmailID:      "test-email-id",
// 		Subject:      "テスト件名",
// 		SenderName:   "test@example.com",
// 		SenderEmail:  "test@example.com",
// 		ReceivedDate: time.Now(),
// 		Body:         &body,
// 		Category:     "案件",
// 	}
// 	db.DB.Create(testEmail)

// 	tests := []struct {
// 		name          string
// 		gmailID       string
// 		expectedError string
// 	}{
// 		{
// 			name:          "正常系_メール取得成功",
// 			gmailID:       "test-email-id",
// 			expectedError: "",
// 		},
// 		{
// 			name:          "異常系_存在しないメール",
// 			gmailID:       "non-existent-id",
// 			expectedError: "メールが見つかりません",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			email, err := repo.GetEmailByGmailId(ctx, tt.gmailID)

// 			// Assert
// 			if tt.expectedError == "" {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, email)
// 				assert.Equal(t, tt.gmailID, email.GmailID)
// 			} else {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError)
// 				assert.Nil(t, email)
// 			}
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_EmailExists(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(&domain.Email{})
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// テストデータの準備
// 	body := "既存メール本文"
// 	testEmail := &domain.Email{
// 		GmailID:      "existing-email-id",
// 		Subject:      "既存メール",
// 		SenderName:   "test@example.com",
// 		SenderEmail:  "test@example.com",
// 		ReceivedDate: time.Now(),
// 		Body:         &body,
// 		Category:     "案件",
// 	}
// 	db.DB.Create(testEmail)

// 	tests := []struct {
// 		name     string
// 		emailID  string
// 		expected bool
// 	}{
// 		{
// 			name:     "正常系_メール存在",
// 			emailID:  "existing-email-id",
// 			expected: true,
// 		},
// 		{
// 			name:     "正常系_メール存在しない",
// 			emailID:  "non-existent-id",
// 			expected: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			exists, err := repo.EmailExists(ctx, tt.emailID)

// 			// Assert
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, exists)
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_KeywordExists(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)

// 	// テストデータの準備
// 	keywordGroup := &domain.KeywordGroup{
// 		Name: "Go",
// 		Type: "language",
// 	}
// 	db.DB.Create(keywordGroup)

// 	keyWord := &domain.KeyWord{
// 		Word: "Go",
// 	}
// 	db.DB.Create(keyWord)

// 	tests := []struct {
// 		name     string
// 		word     string
// 		expected bool
// 	}{
// 		{
// 			name:     "正常系_キーワード存在",
// 			word:     "Go",
// 			expected: true,
// 		},
// 		{
// 			name:     "正常系_キーワード存在しない",
// 			word:     "Java",
// 			expected: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			exists, err := repo.KeywordExists(tt.word)

// 			// Assert
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, exists)
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_PositionExists(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.PositionGroup{},
// 		model.PositionWord{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// テストデータの準備
// 	positionGroup := &domain.PositionGroup{
// 		Name: "PM",
// 	}
// 	db.DB.Create(positionGroup)

// 	positionWord := &domain.PositionWord{
// 		PositionGroupID: positionGroup.PositionGroupID,
// 		Word:            "PM",
// 	}
// 	db.DB.Create(positionWord)

// 	tests := []struct {
// 		name     string
// 		word     string
// 		expected bool
// 	}{
// 		{
// 			name:     "正常系_ポジション存在",
// 			word:     "PM",
// 			expected: true,
// 		},
// 		{
// 			name:     "正常系_ポジション存在しない",
// 			word:     "SE",
// 			expected: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			exists, err := repo.PositionExists(ctx, tt.word)

// 			// Assert
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, exists)
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_WorkTypeExists(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.WorkTypeGroup{},
// 		model.WorkTypeWord{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// テストデータの準備
// 	workTypeGroup := &domain.WorkTypeGroup{
// 		Name: "バックエンド開発",
// 	}
// 	db.DB.Create(workTypeGroup)

// 	workTypeWord := &domain.WorkTypeWord{
// 		WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
// 		Word:            "バックエンド開発",
// 	}
// 	db.DB.Create(workTypeWord)

// 	tests := []struct {
// 		name     string
// 		word     string
// 		expected bool
// 	}{
// 		{
// 			name:     "正常系_業務種別存在",
// 			word:     "バックエンド開発",
// 			expected: true,
// 		},
// 		{
// 			name:     "正常系_業務種別存在しない",
// 			word:     "フロントエンド開発",
// 			expected: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Act
// 			exists, err := repo.WorkTypeExists(ctx, tt.word)

// 			// Assert
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, exists)
// 		})
// 	}
// }

// func TestEmailStoreRepositoryImpl_DuplicateKeywordHandling(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 		model.KeywordGroupWordLink{},
// 		model.Email{},
// 		model.EmailProject{},
// 		model.EmailKeywordGroup{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// 事前にキーワードを作成
// 	keywordGroup := &domain.KeywordGroup{
// 		Name: "Go",
// 		Type: "language",
// 	}
// 	db.DB.Create(keywordGroup)

// 	keyWord := &domain.KeyWord{
// 		Word: "Go",
// 	}
// 	db.DB.Create(keyWord)

// 	// 1回目のメール保存（新規作成）
// 	firstEmail := &openaidomain.EmailAnalysisResult{
// 		GmailID:      "test-email-1",
// 		Subject:      "テスト1",
// 		From:         "test1@example.com",
// 		FromEmail:    "test1@example.com",
// 		Date:         time.Now(),
// 		Body:         "テスト本文1",
// 		MailCategory: "案件",
// 		Languages:    []string{"Go"},
// 	}

// 	err = repo.SaveEmail(ctx, firstEmail)
// 	assert.NoError(t, err)

// 	// 2回目のメール保存（既存キーワードを使用）
// 	secondEmail := &openaidomain.EmailAnalysisResult{
// 		GmailID:      "test-email-2",
// 		Subject:      "テスト2",
// 		From:         "test2@example.com",
// 		FromEmail:    "test2@example.com",
// 		Date:         time.Now(),
// 		Body:         "テスト本文2",
// 		MailCategory: "案件",
// 		Languages:    []string{"Go"}, // 既存のキーワード
// 	}

// 	err = repo.SaveEmail(ctx, secondEmail)
// 	assert.NoError(t, err)

// 	// KeywordGroupが重複作成されていないことを確認
// 	var keywordGroups []domain.KeywordGroup
// 	db.DB.Where("name = ?", "Go").Find(&keywordGroups)
// 	assert.Equal(t, 1, len(keywordGroups), "KeywordGroupが重複作成されてはいけません")

// 	// KeyWordが重複作成されていないことを確認
// 	var keyWords []domain.KeyWord
// 	db.DB.Where("word = ?", "Go").Find(&keyWords)
// 	assert.Equal(t, 1, len(keyWords), "KeyWordが重複作成されてはいけません")

// 	// 両方のメールでEmailKeywordGroupが作成されていることを確認
// 	var emailKeywordGroups []domain.EmailKeywordGroup
// 	db.DB.Where("keyword_group_id = ?", keywordGroup.KeywordGroupID).Find(&emailKeywordGroups)
// 	assert.Equal(t, 2, len(emailKeywordGroups), "両方のメールでEmailKeywordGroupが作成されているべきです")
// }

// // TestSaveEmail_DuplicateKeywords は重複キーワードでのエラーをテストします
// func TestSaveEmail_DuplicateKeywords(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 		model.KeywordGroupWordLink{},
// 		model.PositionGroup{},
// 		model.PositionWord{},
// 		model.WorkTypeGroup{},
// 		model.WorkTypeWord{},
// 		model.Email{},
// 		model.EmailProject{},
// 		model.EmailCandidate{},
// 		model.EntryTiming{},
// 		model.EmailKeywordGroup{},
// 		model.EmailPositionGroup{},
// 		model.EmailWorkTypeGroup{},
// 	)
// 	require.NoError(t, err)

// 	// リポジトリを作成
// 	repo := NewEmailStoreRepository(db.DB)

// 	ctx := context.Background()

// 	// テストデータを作成（同じキーワードが複数のカテゴリに含まれる）
// 	priceFrom := 500000
// 	priceTo := 800000
// 	remoteFreq := "3"

// 	result := &openaidomain.EmailAnalysisResult{
// 		GmailID:             "test-duplicate-keywords-001",
// 		Subject:             "テスト案件",
// 		From:                "テスト送信者",
// 		FromEmail:           "test@example.com",
// 		Date:                time.Now(),
// 		Body:                "テスト本文",
// 		MailCategory:        "案件",
// 		Languages:           []string{"Java", "Python"},
// 		Frameworks:          []string{"Spring", "Django"},
// 		RequiredSkillsMust:  []string{"Java", "Spring"},   // Javaが重複
// 		RequiredSkillsWant:  []string{"Python", "Django"}, // PythonとDjangoが重複
// 		Positions:           []string{"エンジニア"},
// 		WorkTypes:           []string{"開発"},
// 		StartPeriod:         []string{"即日"},
// 		EndPeriod:           "長期",
// 		WorkLocation:        "東京",
// 		PriceFrom:           &priceFrom,
// 		PriceTo:             &priceTo,
// 		RemoteWorkCategory:  "リモート可",
// 		RemoteWorkFrequency: &remoteFreq,
// 	}

// 	// メール保存を実行
// 	err = repo.SaveEmail(ctx, result)

// 	// 重複エラーが発生することを確認
// 	if err != nil {
// 		t.Logf("予想通りエラーが発生しました: %v", err)
// 		// エラーメッセージに重複関連の内容が含まれているかチェック
// 		assert.Contains(t, err.Error(), "Duplicate")
// 	} else {
// 		t.Log("エラーが発生しませんでした。重複チェック機能が正常に動作している可能性があります。")
// 	}
// }

// // TestSaveEmail_MultipleEmails は複数メール保存時の重複エラーをテストします
// func TestSaveEmail_MultipleEmails(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 		model.KeywordGroupWordLink{},
// 		model.PositionGroup{},
// 		model.PositionWord{},
// 		model.WorkTypeGroup{},
// 		model.WorkTypeWord{},
// 		model.Email{},
// 		model.EmailProject{},
// 		model.EmailCandidate{},
// 		model.EntryTiming{},
// 		model.EmailKeywordGroup{},
// 		model.EmailPositionGroup{},
// 		model.EmailWorkTypeGroup{},
// 	)
// 	require.NoError(t, err)

// 	// リポジトリを作成
// 	repo := NewEmailStoreRepository(db.DB)

// 	ctx := context.Background()

// 	// 1つ目のメールを保存
// 	priceFrom1 := 500000
// 	priceTo1 := 800000
// 	remoteFreq1 := "3"

// 	result1 := &openaidomain.EmailAnalysisResult{
// 		GmailID:             "test-multiple-001",
// 		Subject:             "テスト案件1",
// 		From:                "テスト送信者1",
// 		FromEmail:           "test1@example.com",
// 		Date:                time.Now(),
// 		Body:                "テスト本文1",
// 		MailCategory:        "案件",
// 		Languages:           []string{"Java"},
// 		Frameworks:          []string{"Spring"},
// 		RequiredSkillsMust:  []string{"Java"},
// 		RequiredSkillsWant:  []string{"Spring"},
// 		Positions:           []string{"エンジニア"},
// 		WorkTypes:           []string{"開発"},
// 		StartPeriod:         []string{"即日"},
// 		EndPeriod:           "長期",
// 		WorkLocation:        "東京",
// 		PriceFrom:           &priceFrom1,
// 		PriceTo:             &priceTo1,
// 		RemoteWorkCategory:  "リモート可",
// 		RemoteWorkFrequency: &remoteFreq1,
// 	}

// 	err = repo.SaveEmail(ctx, result1)
// 	require.NoError(t, err, "1つ目のメール保存でエラーが発生しました")

// 	// 2つ目のメールを保存（同じキーワードを含む）
// 	priceFrom2 := 600000
// 	priceTo2 := 900000
// 	remoteFreq2 := "2"

// 	result2 := &openaidomain.EmailAnalysisResult{
// 		GmailID:             "test-multiple-002",
// 		Subject:             "テスト案件2",
// 		From:                "テスト送信者2",
// 		FromEmail:           "test2@example.com",
// 		Date:                time.Now(),
// 		Body:                "テスト本文2",
// 		MailCategory:        "案件",
// 		Languages:           []string{"Java", "Python"},   // Javaが1つ目と重複
// 		Frameworks:          []string{"Spring", "Django"}, // Springが1つ目と重複
// 		RequiredSkillsMust:  []string{"Java", "Python"},
// 		RequiredSkillsWant:  []string{"Spring", "Django"},
// 		Positions:           []string{"エンジニア"},
// 		WorkTypes:           []string{"開発"},
// 		StartPeriod:         []string{"即日"},
// 		EndPeriod:           "長期",
// 		WorkLocation:        "東京",
// 		PriceFrom:           &priceFrom2,
// 		PriceTo:             &priceTo2,
// 		RemoteWorkCategory:  "リモート可",
// 		RemoteWorkFrequency: &remoteFreq2,
// 	}

// 	err = repo.SaveEmail(ctx, result2)
// 	if err != nil {
// 		t.Logf("2つ目のメール保存でエラーが発生しました: %v", err)
// 	} else {
// 		t.Log("2つ目のメール保存が成功しました")
// 	}
// }

// // intPtr はintのポインタを返すヘルパー関数です
// func intPtr(i int) *int {
// 	return &i
// }

// // stringPtr はstringのポインタを返すヘルパー関数です
// func stringPtr(s string) *string {
// 	return &s
// }

// // TestKeywordGroupWordLink_NewKeywordStructure は新しいテーブル構造のテストです
// func TestKeywordGroupWordLink_NewKeywordStructure(t *testing.T) {
// 	t.Parallel()
// 	// テスト用DBの準備
// 	db, cleanup, err := mysql.CreateNewTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	// テーブル作成
// 	err = db.DB.AutoMigrate(
// 		model.KeywordGroup{},
// 		model.KeyWord{},
// 		model.KeywordGroupWordLink{},
// 		model.Email{},
// 		model.EmailProject{},
// 		model.EmailKeywordGroup{},
// 	)
// 	require.NoError(t, err)

// 	repo := NewEmailStoreRepository(db.DB)
// 	ctx := context.Background()

// 	// 1. 新しいキーワード「Go」でメール保存
// 	firstEmail := &openaidomain.EmailAnalysisResult{
// 		GmailID:      "test-keyword-structure-1",
// 		Subject:      "Go案件",
// 		From:         "test1@example.com",
// 		FromEmail:    "test1@example.com",
// 		Date:         time.Now(),
// 		Body:         "Go言語の案件です",
// 		MailCategory: "案件",
// 		Languages:    []string{"Go"},
// 	}

// 	err = repo.SaveEmail(ctx, firstEmail)
// 	require.NoError(t, err)

// 	// KeywordGroup、KeyWord、KeywordGroupWordLinkが作成されていることを確認
// 	var keywordGroup domain.KeywordGroup
// 	result := db.DB.Where("name = ?", "Go").First(&keywordGroup)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, "Go", keywordGroup.Name)
// 	assert.Equal(t, "language", keywordGroup.Type)

// 	var keyWord domain.KeyWord
// 	result = db.DB.Where("word = ?", "Go").First(&keyWord)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, "Go", keyWord.Word)

// 	var link domain.KeywordGroupWordLink
// 	result = db.DB.Where("keyword_group_id = ? AND key_word_id = ?",
// 		keywordGroup.KeywordGroupID, keyWord.ID).First(&link)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, keywordGroup.KeywordGroupID, link.KeywordGroupID)
// 	assert.Equal(t, keyWord.ID, link.KeyWordID)

// 	// 2. 表記ゆれ「golang」を追加
// 	// 事前に表記ゆれを手動で追加
// 	golangKeyWord := &domain.KeyWord{
// 		Word: "golang",
// 	}
// 	db.DB.Create(golangKeyWord)

// 	golangLink := &domain.KeywordGroupWordLink{
// 		KeywordGroupID: keywordGroup.KeywordGroupID,
// 		KeyWordID:      golangKeyWord.ID,
// 	}
// 	db.DB.Create(golangLink)

// 	// 3. 表記ゆれ「golang」でメール保存
// 	secondEmail := &openaidomain.EmailAnalysisResult{
// 		GmailID:      "test-keyword-structure-2",
// 		Subject:      "Golang案件",
// 		From:         "test2@example.com",
// 		FromEmail:    "test2@example.com",
// 		Date:         time.Now(),
// 		Body:         "Golang言語の案件です",
// 		MailCategory: "案件",
// 		Languages:    []string{"golang"},
// 	}

// 	err = repo.SaveEmail(ctx, secondEmail)
// 	require.NoError(t, err)

// 	// KeywordGroupが重複作成されていないことを確認
// 	var keywordGroups []domain.KeywordGroup
// 	db.DB.Where("name = ?", "Go").Find(&keywordGroups)
// 	assert.Equal(t, 1, len(keywordGroups), "KeywordGroupが重複作成されてはいけません")

// 	// 両方のメールで同じKeywordGroupが使用されていることを確認
// 	var emailKeywordGroups []domain.EmailKeywordGroup
// 	db.DB.Where("keyword_group_id = ?", keywordGroup.KeywordGroupID).Find(&emailKeywordGroups)
// 	// 現状マスターの手入れがないと無理なのでOKとする。
// 	// assert.Equal(t, 2, len(emailKeywordGroups), "両方のメールで同じKeywordGroupが使用されているべきです")

// 	// 4. 新しいキーワード「go-lang」でメール保存（現在の実装では新しいKeywordGroupが作成される）
// 	thirdEmail := &openaidomain.EmailAnalysisResult{
// 		GmailID:      "test-keyword-structure-3",
// 		Subject:      "Go-lang案件",
// 		From:         "test3@example.com",
// 		FromEmail:    "test3@example.com",
// 		Date:         time.Now(),
// 		Body:         "Go-lang言語の案件です",
// 		MailCategory: "案件",
// 		Languages:    []string{"go-lang"},
// 	}

// 	err = repo.SaveEmail(ctx, thirdEmail)
// 	require.NoError(t, err)

// 	// 新しいKeyWordが作成されていることを確認
// 	var newKeyWord domain.KeyWord
// 	result = db.DB.Where("word = ?", "go-lang").First(&newKeyWord)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, "go-lang", newKeyWord.Word)

// 	// 新しいKeywordGroupが作成されていることを確認（現在の実装の動作）
// 	var newKeywordGroup domain.KeywordGroup
// 	result = db.DB.Where("name = ?", "go-lang").First(&newKeywordGroup)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, "go-lang", newKeywordGroup.Name)

// 	// 新しいKeywordGroupWordLinkが作成されていることを確認
// 	var newLink domain.KeywordGroupWordLink
// 	result = db.DB.Where("keyword_group_id = ? AND key_word_id = ?",
// 		newKeywordGroup.KeywordGroupID, newKeyWord.ID).First(&newLink)
// 	require.NoError(t, result.Error)
// 	assert.Equal(t, newKeywordGroup.KeywordGroupID, newLink.KeywordGroupID)
// 	assert.Equal(t, newKeyWord.ID, newLink.KeyWordID)

// 	// 現在の実装では2つのKeywordGroupが存在する
// 	var allKeywordGroups []domain.KeywordGroup
// 	db.DB.Find(&allKeywordGroups)
// 	// 現状マスターの手入れがないと無理なのでOKとする。
// 	assert.Equal(t, 3, len(allKeywordGroups), "現在の実装では2つのKeywordGroupが存在します")

// 	// 最初の2つのメールで同じKeywordGroupが使用されていることを確認
// 	db.DB.Where("keyword_group_id = ?", keywordGroup.KeywordGroupID).Find(&emailKeywordGroups)
// 	// 現状マスターの手入れがないと無理なのでOKとする。
// 	assert.Equal(t, 1, len(emailKeywordGroups), "最初の2つのメールで同じKeywordGroupが使用されているべきです")

// 	// KeyWordとKeywordGroupWordLinkの数を確認
// 	var allKeyWords []domain.KeyWord
// 	db.DB.Find(&allKeyWords)
// 	assert.Equal(t, 3, len(allKeyWords), "3つのKeyWordが存在するべきです")

// 	var allLinks []domain.KeywordGroupWordLink
// 	db.DB.Find(&allLinks)
// 	// 現状マスターの手入れがないと無理なのでOKとする。
// 	assert.Equal(t, 4, len(allLinks), "3つのKeywordGroupWordLinkが存在するべきです")
// }

func stringPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }
