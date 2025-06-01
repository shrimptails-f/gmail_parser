// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	"business/internal/emailstore/domain"
	openaidomain "business/internal/openai/domain"
	"context"
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
func NewEmailStoreRepository(db *gorm.DB) EmailStoreRepository {
	return &EmailStoreRepositoryImpl{
		db: db,
	}
}

// SaveEmail はメール分析結果をデータベースに保存します
func (r *EmailStoreRepositoryImpl) SaveEmail(ctx context.Context, result *openaidomain.EmailAnalysisResult) error {
	// 重複チェック
	exists, err := r.EmailExists(ctx, result.GmailID)
	if err != nil {
		return fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

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

	// Emailテーブルに保存
	body := result.Body
	email := &domain.Email{
		GmailID:      result.GmailID,
		Subject:      result.Subject,
		SenderName:   result.From,
		SenderEmail:  result.FromEmail,
		ReceivedDate: result.Date,
		Body:         &body,
		Category:     "案件", // デフォルトで案件として設定
	}

	if err := tx.Create(&email).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	// 案件メールの場合、詳細情報を保存
	if result.MailCategory == "案件" {
		if err := r.saveProjectDetails(tx, result, *email); err != nil {
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

// SaveEmailMultiple は複数案件対応のメール分析結果をデータベースに保存します
func (r *EmailStoreRepositoryImpl) SaveEmailMultiple(ctx context.Context, result *openaidomain.EmailAnalysisMultipleResult) error {
	// 重複チェック
	exists, err := r.EmailExists(ctx, result.GmailID)
	if err != nil {
		return fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

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

	// Emailテーブルに保存
	body := result.Body

	// 案件情報を保存
	for _, project := range result.Projects {
		email := &domain.Email{
			GmailID:      result.GmailID,
			Subject:      result.Subject,
			SenderName:   result.From,
			SenderEmail:  result.FromEmail,
			ReceivedDate: result.Date,
			Body:         &body,
			Category:     result.MailCategory,
		}
		if err := tx.Create(&email).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("メール保存エラー: %w", err)
		}

		if err := r.saveProjectDetailsMultiple(tx, email.ID, &project); err != nil {
			tx.Rollback()
			return fmt.Errorf("案件詳細保存エラー: %w", err)
		}
	}

	// 人材情報を保存
	for _, candidate := range result.Candidates {
		email := &domain.Email{
			GmailID:      result.GmailID,
			Subject:      result.Subject,
			SenderName:   result.From,
			SenderEmail:  result.FromEmail,
			ReceivedDate: result.Date,
			Body:         &body,
			Category:     result.MailCategory,
		}
		if err := tx.Create(&email).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("メール保存エラー: %w", err)
		}

		if err := r.saveCandidateDetails(tx, email.ID, &candidate); err != nil {
			tx.Rollback()
			return fmt.Errorf("人材詳細保存エラー: %w", err)
		}
	}

	// トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	return nil
}

// saveProjectDetails は案件メールの詳細情報を保存します
func (r *EmailStoreRepositoryImpl) saveProjectDetails(tx *gorm.DB, result *openaidomain.EmailAnalysisResult, email domain.Email) error {
	// EmailProjectを保存
	workLocation := result.WorkLocation
	endTiming := result.EndPeriod
	remoteType := result.RemoteWorkCategory
	languages := strings.Join(result.Languages, ",")
	frameworks := strings.Join(result.Frameworks, ",")
	positions := strings.Join(result.Positions, ",")
	workTypes := strings.Join(result.WorkTypes, ",")
	mustSkills := strings.Join(result.RequiredSkillsMust, ",")
	wantSkills := strings.Join(result.RequiredSkillsWant, ",")

	emailProject := &domain.EmailProject{
		EmailID:         email.ID,
		WorkLocation:    &workLocation,
		EndTiming:       &endTiming,
		PriceFrom:       result.PriceFrom,
		PriceTo:         result.PriceTo,
		RemoteType:      &remoteType,
		RemoteFrequency: result.RemoteWorkFrequency,
		Languages:       &languages,
		Frameworks:      &frameworks,
		Positions:       &positions,
		WorkTypes:       &workTypes,
		MustSkills:      &mustSkills,
		WantSkills:      &wantSkills,
	}

	if err := tx.Create(emailProject).Error; err != nil {
		return fmt.Errorf("EmailProject保存エラー: %w", err)
	}

	// EntryTimingを保存
	if err := r.saveEntryTimings(tx, emailProject.EmailID, result.StartPeriod); err != nil {
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

// saveProjectDetailsMultiple は複数案件対応の案件詳細情報を保存します
func (r *EmailStoreRepositoryImpl) saveProjectDetailsMultiple(tx *gorm.DB, emailID uint, project *openaidomain.ProjectAnalysisResult) error {
	// EmailProjectを保存
	workLocation := project.WorkLocation
	endTiming := project.EndPeriod
	remoteType := project.RemoteWorkCategory
	languages := strings.Join(project.Languages, ",")
	frameworks := strings.Join(project.Frameworks, ",")
	positions := strings.Join(project.Positions, ",")
	workTypes := strings.Join(project.WorkTypes, ",")
	mustSkills := strings.Join(project.RequiredSkillsMust, ",")
	wantSkills := strings.Join(project.RequiredSkillsWant, ",")
	projectTitle := project.ProjectName

	emailProject := &domain.EmailProject{
		EmailID:         emailID,
		ProjectTitle:    &projectTitle,
		WorkLocation:    &workLocation,
		EndTiming:       &endTiming,
		PriceFrom:       project.PriceFrom,
		PriceTo:         project.PriceTo,
		RemoteType:      &remoteType,
		RemoteFrequency: project.RemoteWorkFrequency,
		Languages:       &languages,
		Frameworks:      &frameworks,
		Positions:       &positions,
		WorkTypes:       &workTypes,
		MustSkills:      &mustSkills,
		WantSkills:      &wantSkills,
	}

	if err := tx.Create(emailProject).Error; err != nil {
		return fmt.Errorf("EmailProject保存エラー: %w", err)
	}

	// EntryTimingを保存
	if err := r.saveEntryTimings(tx, emailProject.EmailID, project.StartPeriod); err != nil {
		return fmt.Errorf("EntryTiming保存エラー: %w", err)
	}

	// キーワード関連を保存
	if err := r.saveKeywordsByTypeMultiple(tx, emailID, project.Languages, "LANGUAGE"); err != nil {
		return fmt.Errorf("言語キーワード保存エラー: %w", err)
	}
	if err := r.saveKeywordsByTypeMultiple(tx, emailID, project.Frameworks, "FRAMEWORK"); err != nil {
		return fmt.Errorf("フレームワークキーワード保存エラー: %w", err)
	}
	if err := r.saveKeywordsByTypeMultiple(tx, emailID, project.RequiredSkillsMust, "MUST"); err != nil {
		return fmt.Errorf("必須スキル保存エラー: %w", err)
	}
	if err := r.saveKeywordsByTypeMultiple(tx, emailID, project.RequiredSkillsWant, "WANT"); err != nil {
		return fmt.Errorf("希望スキル保存エラー: %w", err)
	}

	// ポジション関連を保存
	if err := r.savePositionsMultiple(tx, emailID, project.Positions); err != nil {
		return fmt.Errorf("ポジション保存エラー: %w", err)
	}

	// 業務種別関連を保存
	if err := r.saveWorkTypesMultiple(tx, emailID, project.WorkTypes); err != nil {
		return fmt.Errorf("業務種別保存エラー: %w", err)
	}

	return nil
}

// saveCandidateDetails は人材詳細情報を保存します
func (r *EmailStoreRepositoryImpl) saveCandidateDetails(tx *gorm.DB, emailID uint, candidate *openaidomain.CandidateAnalysisResult) error {
	// EmailCandidateを保存
	emailCandidate := &domain.EmailCandidate{
		EmailID: emailID,
	}

	if err := tx.Create(emailCandidate).Error; err != nil {
		return fmt.Errorf("EmailCandidate保存エラー: %w", err)
	}

	// TODO: 人材情報の詳細フィールドが追加されたら実装
	// 現在はEmailCandidateテーブルの基本情報のみ保存

	return nil
}

// saveKeywordsByTypeMultiple は複数案件対応のキーワード保存です
func (r *EmailStoreRepositoryImpl) saveKeywordsByTypeMultiple(tx *gorm.DB, emailID uint, keywords []string, keywordType string) error {
	return r.saveKeywordsByType(tx, emailID, keywords, keywordType)
}

// savePositionsMultiple は複数案件対応のポジション保存です
func (r *EmailStoreRepositoryImpl) savePositionsMultiple(tx *gorm.DB, emailID uint, positions []string) error {
	for _, position := range positions {
		if position == "" {
			continue
		}

		// PositionGroupを取得または作成
		positionGroup, err := r.getOrCreatePositionGroup(tx, position)
		if err != nil {
			return fmt.Errorf("PositionGroup取得/作成エラー: %w", err)
		}

		// EmailPositionGroupを作成
		emailPositionGroup := &domain.EmailPositionGroup{
			EmailID:         emailID,
			PositionGroupID: positionGroup.PositionGroupID,
		}

		if err := tx.Create(emailPositionGroup).Error; err != nil {
			return fmt.Errorf("EmailPositionGroup保存エラー: %w", err)
		}
	}
	return nil
}

// saveWorkTypesMultiple は複数案件対応の業務種別保存です
func (r *EmailStoreRepositoryImpl) saveWorkTypesMultiple(tx *gorm.DB, emailID uint, workTypes []string) error {
	for _, workType := range workTypes {
		if workType == "" {
			continue
		}

		// WorkTypeGroupを取得または作成
		workTypeGroup, err := r.getOrCreateWorkTypeGroup(tx, workType)
		if err != nil {
			return fmt.Errorf("WorkTypeGroup取得/作成エラー: %w", err)
		}

		// EmailWorkTypeGroupを作成
		emailWorkTypeGroup := &domain.EmailWorkTypeGroup{
			EmailID:         emailID,
			WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
		}

		if err := tx.Create(emailWorkTypeGroup).Error; err != nil {
			return fmt.Errorf("EmailWorkTypeGroup保存エラー: %w", err)
		}
	}
	return nil
}

// saveEntryTimings は入場時期を保存します
func (r *EmailStoreRepositoryImpl) saveEntryTimings(tx *gorm.DB, emailProjectID uint, startPeriods []string) error {
	for _, period := range startPeriods {
		entryTiming := &domain.EntryTiming{
			EmailProjectID: emailProjectID,
			StartDate:      period,
		}
		if err := tx.Create(entryTiming).Error; err != nil {
			return fmt.Errorf("EntryTiming保存エラー: %w", err)
		}
	}
	return nil
}

// saveKeywords はキーワード関連のデータを保存します
func (r *EmailStoreRepositoryImpl) saveKeywords(tx *gorm.DB, result *openaidomain.EmailAnalysisResult, emailId uint) error {
	// 言語
	if err := r.saveKeywordsByType(tx, emailId, result.Languages, "LANGUAGE"); err != nil {
		return err
	}

	// フレームワーク
	if err := r.saveKeywordsByType(tx, emailId, result.Frameworks, "FRAMEWORK"); err != nil {
		return err
	}

	// 必須スキル
	if err := r.saveKeywordsByType(tx, emailId, result.RequiredSkillsMust, "MUST"); err != nil {
		return err
	}

	// 希望スキル
	if err := r.saveKeywordsByType(tx, emailId, result.RequiredSkillsWant, "WANT"); err != nil {
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
		keywordGroup, err := r.getOrCreateKeywordGroup(tx, keyword)
		if err != nil {
			return fmt.Errorf("KeywordGroup取得/作成エラー: %w", err)
		}

		// EmailKeywordGroupを作成
		emailKeywordGroup := &domain.EmailKeywordGroup{
			EmailID:        emailID,
			KeywordGroupID: keywordGroup.KeywordGroupID,
			Type:           keywordType,
		}

		if err := tx.Create(emailKeywordGroup).Error; err != nil {
			return fmt.Errorf("EmailKeywordGroup保存エラー: %w", err)
		}
	}
	return nil
}

func (r *EmailStoreRepositoryImpl) getOrCreateKeywordGroup(tx *gorm.DB, name string) (*domain.KeywordGroup, error) {
	var keywordGroup domain.KeywordGroup
	var keyWord domain.KeyWord

	// 単語を先に取得（存在確認）
	err := tx.Where("word = ?", name).First(&keyWord).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("KeyWord検索エラー: %w", err)
	}

	// グループ存在確認
	err = tx.Where("name = ?", name).First(&keywordGroup).Error
	if err == nil {
		// 存在する場合 → Linkの存在確認
		if keyWord.ID == 0 {
			// 単語がなければ作成
			keyWord = domain.KeyWord{Word: name}
			if err := tx.Create(&keyWord).Error; err != nil {
				return nil, fmt.Errorf("KeyWord作成エラー: %w", err)
			}
		}

		// 中間テーブルが存在するかチェック
		var count int64
		err = tx.Model(&domain.KeywordGroupWordLink{}).
			Where("keyword_group_id = ? AND key_word_id = ?", keywordGroup.KeywordGroupID, keyWord.ID).
			Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("KeywordGroupWordLink確認エラー: %w", err)
		}

		if count == 0 {
			link := &domain.KeywordGroupWordLink{
				KeywordGroupID: keywordGroup.KeywordGroupID,
				KeyWordID:      keyWord.ID,
			}
			if err := tx.Create(link).Error; err != nil {
				return nil, fmt.Errorf("KeywordGroupWordLink作成エラー: %w", err)
			}
		}

		return &keywordGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("KeywordGroup検索エラー: %w", err)
	}

	// グループも単語も存在しない → 新規作成
	keywordGroup = domain.KeywordGroup{
		Name: name,
		Type: "other",
	}
	if err := tx.Create(&keywordGroup).Error; err != nil {
		return nil, fmt.Errorf("KeywordGroup作成エラー: %w", err)
	}

	if keyWord.ID == 0 {
		keyWord = domain.KeyWord{Word: name}
		if err := tx.Create(&keyWord).Error; err != nil {
			return nil, fmt.Errorf("KeyWord作成エラー: %w", err)
		}
	}

	link := &domain.KeywordGroupWordLink{
		KeywordGroupID: keywordGroup.KeywordGroupID,
		KeyWordID:      keyWord.ID,
	}
	if err := tx.Create(link).Error; err != nil {
		return nil, fmt.Errorf("KeywordGroupWordLink作成エラー: %w", err)
	}

	return &keywordGroup, nil
}

// savePositions はポジション関連のデータを保存します
func (r *EmailStoreRepositoryImpl) savePositions(tx *gorm.DB, result *openaidomain.EmailAnalysisResult, emailId uint) error {
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
		emailPositionGroup := &domain.EmailPositionGroup{
			EmailID:         emailId,
			PositionGroupID: positionGroup.PositionGroupID,
		}

		if err := tx.Create(emailPositionGroup).Error; err != nil {
			return fmt.Errorf("EmailPositionGroup保存エラー: %w", err)
		}
	}
	return nil
}

// getOrCreatePositionGroup はPositionGroupを取得または作成します
func (r *EmailStoreRepositoryImpl) getOrCreatePositionGroup(tx *gorm.DB, name string) (*domain.PositionGroup, error) {
	var positionGroup domain.PositionGroup

	// 既存のPositionGroupを検索
	err := tx.Where("name = ?", name).First(&positionGroup).Error
	if err == nil {
		return &positionGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("PositionGroup検索エラー: %w", err)
	}

	// 表記ゆれとして既に存在するかチェック
	var existingPositionWord domain.PositionWord
	err = tx.Where("word = ?", name).First(&existingPositionWord).Error
	if err == nil {
		// 既存の表記ゆれが見つかった場合、対応するPositionGroupを取得
		err = tx.Where("position_group_id = ?", existingPositionWord.PositionGroupID).First(&positionGroup).Error
		if err != nil {
			return nil, fmt.Errorf("既存PositionGroup取得エラー: %w", err)
		}
		return &positionGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("PositionWord検索エラー: %w", err)
	}

	// 新規作成
	positionGroup = domain.PositionGroup{
		Name: name,
	}

	if err := tx.Create(&positionGroup).Error; err != nil {
		return nil, fmt.Errorf("PositionGroup作成エラー: %w", err)
	}

	// PositionWordも作成（表記ゆれとして同じ名前を登録）
	positionWord := &domain.PositionWord{
		PositionGroupID: positionGroup.PositionGroupID,
		Word:            name,
	}

	if err := tx.Create(positionWord).Error; err != nil {
		return nil, fmt.Errorf("PositionWord作成エラー: %w", err)
	}

	return &positionGroup, nil
}

// saveWorkTypes は業務種別関連のデータを保存します
func (r *EmailStoreRepositoryImpl) saveWorkTypes(tx *gorm.DB, result *openaidomain.EmailAnalysisResult, emailId uint) error {
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
		emailWorkTypeGroup := &domain.EmailWorkTypeGroup{
			EmailID:         emailId,
			WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
		}

		if err := tx.Create(emailWorkTypeGroup).Error; err != nil {
			return fmt.Errorf("EmailWorkTypeGroup保存エラー: %w", err)
		}
	}
	return nil
}

// getOrCreateWorkTypeGroup はWorkTypeGroupを取得または作成します
func (r *EmailStoreRepositoryImpl) getOrCreateWorkTypeGroup(tx *gorm.DB, name string) (*domain.WorkTypeGroup, error) {
	var workTypeGroup domain.WorkTypeGroup

	// 既存のWorkTypeGroupを検索
	err := tx.Where("name = ?", name).First(&workTypeGroup).Error
	if err == nil {
		return &workTypeGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("WorkTypeGroup検索エラー: %w", err)
	}

	// 表記ゆれとして既に存在するかチェック
	var existingWorkTypeWord domain.WorkTypeWord
	err = tx.Where("word = ?", name).First(&existingWorkTypeWord).Error
	if err == nil {
		// 既存の表記ゆれが見つかった場合、対応するWorkTypeGroupを取得
		err = tx.Where("work_type_group_id = ?", existingWorkTypeWord.WorkTypeGroupID).First(&workTypeGroup).Error
		if err != nil {
			return nil, fmt.Errorf("既存WorkTypeGroup取得エラー: %w", err)
		}
		return &workTypeGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("WorkTypeWord検索エラー: %w", err)
	}

	// 新規作成
	workTypeGroup = domain.WorkTypeGroup{
		Name: name,
	}

	if err := tx.Create(&workTypeGroup).Error; err != nil {
		return nil, fmt.Errorf("WorkTypeGroup作成エラー: %w", err)
	}

	// WorkTypeWordも作成（表記ゆれとして同じ名前を登録）
	workTypeWord := &domain.WorkTypeWord{
		WorkTypeGroupID: workTypeGroup.WorkTypeGroupID,
		Word:            name,
	}

	if err := tx.Create(workTypeWord).Error; err != nil {
		return nil, fmt.Errorf("WorkTypeWord作成エラー: %w", err)
	}

	return &workTypeGroup, nil
}

// GetEmailByGmailId はIDでメールを取得します
func (r *EmailStoreRepositoryImpl) GetEmailByGmailId(ctx context.Context, gmail_id string) (*domain.Email, error) {
	var email domain.Email
	err := r.db.Where("gmail_id = ?", gmail_id).First(&email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmailNotFound
		}
		return nil, fmt.Errorf("メール取得エラー: %w", err)
	}
	return &email, nil
}

// EmailExists はメールが既に存在するかチェックします
func (r *EmailStoreRepositoryImpl) EmailExists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Email{}).Where("gmail_id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	return count > 0, nil
}

// KeywordExists はキーワードが既に存在するかチェックします
func (r *EmailStoreRepositoryImpl) KeywordExists(word string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.KeyWord{}).Where("word = ?", word).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("キーワード存在チェックエラー: %w", err)
	}
	return count > 0, nil
}

// PositionExists はポジションが既に存在するかチェックします
func (r *EmailStoreRepositoryImpl) PositionExists(ctx context.Context, word string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.PositionWord{}).Where("word = ?", word).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("ポジション存在チェックエラー: %w", err)
	}
	return count > 0, nil
}

// WorkTypeExists は業務種別が既に存在するかチェックします
func (r *EmailStoreRepositoryImpl) WorkTypeExists(ctx context.Context, word string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.WorkTypeWord{}).Where("word = ?", word).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("業務種別存在チェックエラー: %w", err)
	}
	return count > 0, nil
}
