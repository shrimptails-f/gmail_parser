// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	cd "business/internal/common/domain"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// EmailStoreRepositoryImpl はメール保存のリポジトリ実装です
type EmailStoreRepositoryImpl struct {
	db *gorm.DB
}

// NewEmailStoreRepository はメール保存リポジトリを作成します
func NewEmailStoreRepository(db *gorm.DB) *EmailStoreRepositoryImpl {
	return &EmailStoreRepositoryImpl{
		db: db,
	}
}

func (r *EmailStoreRepositoryImpl) SaveEmail(result cd.Email) error {
	// トランザクション開始
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	email := r.setEmail(result)

	if err := tx.Create(&email).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	// 案件メールの場合、詳細情報を保存
	if result.Category == "案件" {
		if err := r.saveProjectDetails(tx, result, email); err != nil {
			tx.Rollback()
			return fmt.Errorf("案件詳細保存エラー: %w", err)
		}
	}

	// トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	return nil
}

func (r *EmailStoreRepositoryImpl) setEmail(result cd.Email) Email {
	return Email{
		GmailID:      result.GmailID,
		Subject:      result.Subject,
		SenderName:   result.SenderName(),
		SenderEmail:  result.SenderEmail(),
		ReceivedDate: result.ReceivedDate,
		Body:         &result.Body,
		Category:     result.Category,
	}
}

// GetEmailByGmailId はIDでメールを取得します
func (r *EmailStoreRepositoryImpl) GetEmailByGmailId(gmail_id string) (Email, error) {
	var email Email
	err := r.db.Where("gmail_id = ?", gmail_id).First(&email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Email{}, ErrEmailNotFound
		}
		return Email{}, fmt.Errorf("メール取得エラー: %w", err)
	}
	return email, nil
}

// EmailExists はメールが既に存在するかチェックします
func (r *EmailStoreRepositoryImpl) EmailExists(id string) (bool, error) {
	var count int64
	err := r.db.Model(Email{}).Where("gmail_id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	return count > 0, nil
}

// saveProjectDetails は案件メールの詳細情報を保存します
func (r *EmailStoreRepositoryImpl) saveProjectDetails(tx *gorm.DB, result cd.Email, email Email) error {
	// EmailProjectを保存
	entryTimings := strings.Join(result.StartPeriod, ",")
	languages := strings.Join(result.Languages, ",")
	frameworks := strings.Join(result.Frameworks, ",")
	positions := strings.Join(result.Positions, ",")
	workTypes := strings.Join(result.WorkTypes, ",")
	mustSkills := strings.Join(result.RequiredSkillsMust, ",")
	wantSkills := strings.Join(result.RequiredSkillsWant, ",")

	emailProject := EmailProject{
		EmailID:         email.ID,
		ProjectTitle:    &result.Summary,
		EntryTiming:     &entryTimings,
		WorkLocation:    &result.WorkLocation,
		EndTiming:       &result.EndPeriod,
		PriceFrom:       result.PriceFrom,
		PriceTo:         result.PriceTo,
		RemoteType:      result.RemoteWorkCategory,
		RemoteFrequency: result.RemoteWorkFrequency,
		Languages:       &languages,
		Frameworks:      &frameworks,
		Positions:       &positions,
		WorkTypes:       &workTypes,
		MustSkills:      &mustSkills,
		WantSkills:      &wantSkills,
	}

	if err := tx.Create(&emailProject).Error; err != nil {
		return fmt.Errorf("EmailProject保存エラー: %w", err)
	}

	// EntryTimingを保存
	if err := r.saveEntryTimings(tx, email.ID, result.StartPeriod); err != nil {
		return fmt.Errorf("EntryTiming保存エラー: %w", err)
	}

	// キーワード関連を保存
	if err := r.saveKeywords(tx, result, email.ID); err != nil {
		return fmt.Errorf("キーワード保存エラー: %w", err)
	}

	// ポジション関連を保存
	if err := r.savePositions(tx, result, email.ID); err != nil {
		return fmt.Errorf("ポジション保存エラー: %w", err)
	}

	// 業務種別関連を保存
	if err := r.saveWorkTypes(tx, result, email.ID); err != nil {
		return fmt.Errorf("業務種別保存エラー: %w", err)
	}

	return nil
}

// saveEntryTimings は入場時期を保存します
func (r *EmailStoreRepositoryImpl) saveEntryTimings(tx *gorm.DB, emailId uint, startPeriods []string) error {
	for _, period := range startPeriods {
		entryTiming := EntryTiming{
			EmailID:   emailId,
			StartDate: period,
		}
		if err := tx.Create(&entryTiming).Error; err != nil {
			return fmt.Errorf("EntryTiming保存エラー: %w", err)
		}
	}
	return nil
}

// saveKeywords はキーワード関連のデータを保存します
func (r *EmailStoreRepositoryImpl) saveKeywords(tx *gorm.DB, result cd.Email, emailId uint) error {
	// 言語
	if err := r.saveKeywordsByType(tx, emailId, result.Languages, "language"); err != nil {
		return err
	}

	// フレームワーク
	if err := r.saveKeywordsByType(tx, emailId, result.Frameworks, "framework"); err != nil {
		return err
	}

	// 必須スキル
	if err := r.saveKeywordsByType(tx, emailId, result.RequiredSkillsMust, "must"); err != nil {
		return err
	}

	// 希望スキル
	if err := r.saveKeywordsByType(tx, emailId, result.RequiredSkillsWant, "want"); err != nil {
		return err
	}

	return nil
}

// saveKeywordsByType は指定されたタイプのキーワードを保存します
func (r *EmailStoreRepositoryImpl) saveKeywordsByType(tx *gorm.DB, emailID uint, keywords []string, keywordType string) error {
	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}

		// KeywordGroupを取得または作成
		keywordGroup, err := r.getOrCreateKeywordGroup(tx, keyword, keywordType)
		if err != nil {
			return fmt.Errorf("KeywordGroup取得/作成エラー: %w", err)
		}

		// EmailKeywordGroupを作成
		emailKeywordGroup := EmailKeywordGroup{
			EmailID:        emailID,
			KeywordGroupID: keywordGroup.KeywordGroupID,
		}

		if err := tx.Create(&emailKeywordGroup).Error; err != nil {
			return fmt.Errorf("EmailKeywordGroup保存エラー: %w", err)
		}
	}
	return nil
}

func (r *EmailStoreRepositoryImpl) getOrCreateKeywordGroup(tx *gorm.DB, name string, keywordType string) (KeywordGroup, error) {
	var keywordGroup KeywordGroup
	var keyWord KeyWord

	// 単語を先に取得（存在確認）
	err := tx.Where("word = ?", name).First(&keyWord).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return KeywordGroup{}, fmt.Errorf("KeyWord検索エラー: %w", err)
	}

	// グループ存在確認
	err = tx.Where("name = ?", name).First(&keywordGroup).Error
	if err == nil {
		// 存在する場合 → Linkの存在確認
		if keyWord.ID == 0 {
			// 単語がなければ作成
			keyWord = KeyWord{Word: name}
			if err := tx.Create(&keyWord).Error; err != nil {
				return KeywordGroup{}, fmt.Errorf("KeyWord作成エラー: %w", err)
			}
		}

		// 中間テーブルが存在するかチェック
		var count int64
		err = tx.Model(KeywordGroupWordLink{}).
			Where("keyword_group_id = ? AND key_word_id = ?", keywordGroup.KeywordGroupID, keyWord.ID).
			Count(&count).Error
		if err != nil {
			return KeywordGroup{}, fmt.Errorf("KeywordGroupWordLink確認エラー: %w", err)
		}

		if count == 0 {
			link := KeywordGroupWordLink{
				KeywordGroupID: keywordGroup.KeywordGroupID,
				KeyWordID:      keyWord.ID,
			}
			if err := tx.Create(&link).Error; err != nil {
				return KeywordGroup{}, fmt.Errorf("KeywordGroupWordLink作成エラー: %w", err)
			}
		}

		return keywordGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return KeywordGroup{}, fmt.Errorf("KeywordGroup検索エラー: %w", err)
	}

	// グループも単語も存在しない → 新規作成
	keywordGroup = KeywordGroup{
		Name: name,
		Type: keywordType,
	}
	if err := tx.Create(&keywordGroup).Error; err != nil {
		return KeywordGroup{}, fmt.Errorf("KeywordGroup作成エラー: %w", err)
	}

	if keyWord.ID == 0 {
		keyWord = KeyWord{Word: name}
		if err := tx.Create(&keyWord).Error; err != nil {
			return KeywordGroup{}, fmt.Errorf("KeyWord作成エラー: %w", err)
		}
	}

	link := KeywordGroupWordLink{
		KeywordGroupID: keywordGroup.KeywordGroupID,
		KeyWordID:      keyWord.ID,
	}
	if err := tx.Create(&link).Error; err != nil {
		return KeywordGroup{}, fmt.Errorf("KeywordGroupWordLink作成エラー: %w", err)
	}

	return keywordGroup, nil
}

// savePositions はポジション関連のデータを保存します
func (r *EmailStoreRepositoryImpl) savePositions(tx *gorm.DB, result cd.Email, emailId uint) error {
	for _, position := range result.Positions {
		if position == "" {
			continue
		}

		// PositionGroupを取得または作成
		positionGroup, err := r.getOrCreatePositionGroup(tx, position)
		if err != nil {
			return fmt.Errorf("PositionGroup取得/作成エラー: %w", err)
		}

		// EmailPositionGroupを作成
		emailPositionGroup := EmailPositionGroup{
			EmailID:         emailId,
			PositionGroupID: positionGroup.PositionGroupID,
		}

		if err := tx.Create(&emailPositionGroup).Error; err != nil {
			return fmt.Errorf("EmailPositionGroup保存エラー: %w", err)
		}
	}
	return nil
}

// getOrCreatePositionGroup はPositionGroupを取得または作成します
func (r *EmailStoreRepositoryImpl) getOrCreatePositionGroup(tx *gorm.DB, name string) (PositionGroup, error) {
	var positionGroup PositionGroup

	// 既存のPositionGroupを検索
	err := tx.Where("name = ?", name).First(&positionGroup).Error
	if err == nil {
		return positionGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return PositionGroup{}, fmt.Errorf("PositionGroup検索エラー: %w", err)
	}

	// 表記ゆれとして既に存在するかチェック
	var existingPositionWord PositionWord
	err = tx.Where("word = ?", name).First(&existingPositionWord).Error
	if err == nil {
		// 既存の表記ゆれが見つかった場合、対応するPositionGroupを取得
		err = tx.Where("position_group_id = ?", existingPositionWord.PositionGroupID).First(&positionGroup).Error
		if err != nil {
			return PositionGroup{}, fmt.Errorf("既存PositionGroup取得エラー: %w", err)
		}
		return positionGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return PositionGroup{}, fmt.Errorf("PositionWord検索エラー: %w", err)
	}

	// 新規作成
	positionGroup = PositionGroup{
		Name: name,
	}

	if err := tx.Create(&positionGroup).Error; err != nil {
		return PositionGroup{}, fmt.Errorf("PositionGroup作成エラー: %w", err)
	}

	// PositionWordも作成（表記ゆれとして同じ名前を登録）
	positionWord := PositionWord{
		PositionGroupID: positionGroup.PositionGroupID,
		Word:            name,
	}

	if err := tx.Create(&positionWord).Error; err != nil {
		return PositionGroup{}, fmt.Errorf("PositionWord作成エラー: %w", err)
	}

	return positionGroup, nil
}

// saveWorkTypes は業務種別関連のデータを保存します
func (r *EmailStoreRepositoryImpl) saveWorkTypes(tx *gorm.DB, result cd.Email, emailId uint) error {
	for _, workType := range result.WorkTypes {
		if workType == "" {
			continue
		}

		// WorkTypeGroupを取得または作成
		workTypeGroup, err := r.getOrCreateWorkTypeGroup(tx, workType)
		if err != nil {
			return fmt.Errorf("WorkTypeGroup取得/作成エラー: %w", err)
		}

		// EmailWorkTypeGroupを作成
		emailWorkTypeGroup := EmailWorkTypeGroup{
			EmailID:         emailId,
			WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
		}

		if err := tx.Create(&emailWorkTypeGroup).Error; err != nil {
			return fmt.Errorf("EmailWorkTypeGroup保存エラー: %w", err)
		}
	}
	return nil
}

// getOrCreateWorkTypeGroup はWorkTypeGroupを取得または作成します
func (r *EmailStoreRepositoryImpl) getOrCreateWorkTypeGroup(tx *gorm.DB, name string) (WorkTypeGroup, error) {
	var workTypeGroup WorkTypeGroup

	// 既存のWorkTypeGroupを検索
	err := tx.Where("name = ?", name).First(&workTypeGroup).Error
	if err == nil {
		return workTypeGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return WorkTypeGroup{}, fmt.Errorf("WorkTypeGroup検索エラー: %w", err)
	}

	// 表記ゆれとして既に存在するかチェック
	var existingWorkTypeWord WorkTypeWord
	err = tx.Where("word = ?", name).First(&existingWorkTypeWord).Error
	if err == nil {
		// 既存の表記ゆれが見つかった場合、対応するWorkTypeGroupを取得
		err = tx.Where("work_type_group_id = ?", existingWorkTypeWord.WorkTypeGroupID).First(&workTypeGroup).Error
		if err != nil {
			return WorkTypeGroup{}, fmt.Errorf("既存WorkTypeGroup取得エラー: %w", err)
		}
		return workTypeGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return WorkTypeGroup{}, fmt.Errorf("WorkTypeWord検索エラー: %w", err)
	}

	// 新規作成
	workTypeGroup = WorkTypeGroup{
		Name: name,
	}

	if err := tx.Create(&workTypeGroup).Error; err != nil {
		return WorkTypeGroup{}, fmt.Errorf("WorkTypeGroup作成エラー: %w", err)
	}

	// WorkTypeWordも作成（表記ゆれとして同じ名前を登録）
	workTypeWord := WorkTypeWord{
		WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
		Word:            name,
	}

	if err := tx.Create(&workTypeWord).Error; err != nil {
		return WorkTypeGroup{}, fmt.Errorf("WorkTypeWord作成エラー: %w", err)
	}

	return workTypeGroup, nil
}
