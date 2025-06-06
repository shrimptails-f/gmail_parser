// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	"business/internal/emailstore/domain"
	"errors"
	"fmt"

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

func (r *EmailStoreRepositoryImpl) SaveEmail(email Email) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ---------- ネストが深い構造体を作成する。----------
	// KeywordGroup保存
	for i := range email.EmailKeywordGroups {
		kg := &email.EmailKeywordGroups[i].KeywordGroup
		if err := tx.Create(kg).Error; err != nil {
			tx.Rollback()
			return err
		}
		email.EmailKeywordGroups[i].KeywordGroupID = kg.KeywordGroupID
	}

	// PositionGroup保存
	for i := range email.EmailPositionGroups {
		pg := &email.EmailPositionGroups[i].PositionGroup
		if err := tx.Create(pg).Error; err != nil {
			tx.Rollback()
			return err
		}
		email.EmailPositionGroups[i].PositionGroupID = pg.PositionGroupID
	}

	// WorkTypeGroup保存
	for i := range email.EmailWorkTypeGroups {
		wtg := &email.EmailWorkTypeGroups[i].WorkTypeGroup
		if err := tx.Create(wtg).Error; err != nil {
			tx.Rollback()
			return err
		}
		email.EmailWorkTypeGroups[i].WorkTypeGroupID = wtg.WorkTypeGroupID
	}
	// ------------------------------------------------------------------------

	// Email保存（最後）
	if err := tx.Create(&email).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	return nil
}

// GetEmailByGmailId はIDでメールを取得します
func (r *EmailStoreRepositoryImpl) GetEmailByGmailId(gmail_id string) (*domain.Email, error) {
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
func (r *EmailStoreRepositoryImpl) EmailExists(id string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Email{}).Where("gmail_id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	return count > 0, nil
}

// GetkeywordGroups はキーワードグループを name で一括取得します
func (r *EmailStoreRepositoryImpl) GetkeywordGroups(names []string) ([]KeywordGroup, error) {
	var groups []KeywordGroup
	if err := r.db.Where("name IN ?", names).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetKeywords はキーワードを word で一括取得します
func (r *EmailStoreRepositoryImpl) GetKeywords(words []string) ([]KeyWord, error) {
	var keywords []KeyWord
	if err := r.db.Where("word IN ?", words).Find(&keywords).Error; err != nil {
		return nil, err
	}
	return keywords, nil
}

// GetPositionGroups はポジショングループを name で一括取得します
func (r *EmailStoreRepositoryImpl) GetPositionGroups(names []string) ([]PositionGroup, error) {
	var groups []PositionGroup
	if err := r.db.Where("name IN ?", names).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetPositionWords はポジションの表記ゆれを word で一括取得します
func (r *EmailStoreRepositoryImpl) GetPositionWords(words []string) ([]PositionWord, error) {
	var positionWords []PositionWord
	if err := r.db.Where("word IN ?", words).Find(&positionWords).Error; err != nil {
		return nil, err
	}
	return positionWords, nil
}

// GetWorkTypeWords は業務種別の表記ゆれを word で一括取得します
func (r *EmailStoreRepositoryImpl) GetWorkTypeWords(words []string) ([]WorkTypeWord, error) {
	var workTypeWords []WorkTypeWord
	if err := r.db.Where("word IN ?", words).Find(&workTypeWords).Error; err != nil {
		return nil, err
	}
	return workTypeWords, nil
}

// GetWorkTypeGroups は業務種別のグループを name で一括取得します
func (r *EmailStoreRepositoryImpl) GetWorkTypeGroups(names []string) ([]WorkTypeGroup, error) {
	var groups []WorkTypeGroup
	if err := r.db.Where("name IN ?", names).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
