// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	"business/internal/emailstore/application"
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
func NewEmailStoreRepository(db *gorm.DB) application.EmailStoreRepository {
	return &EmailStoreRepositoryImpl{
		db: db,
	}
}

// SaveEmail はメール分析結果をデータベースに保存します
func (r *EmailStoreRepositoryImpl) SaveEmail(ctx context.Context, result *openaidomain.EmailAnalysisResult) error {
	// 重複チェック
	exists, err := r.EmailExists(ctx, result.ID)
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
	email := &domain.Email{
		ID:        result.ID,
		Subject:   result.Subject,
		From:      result.From,
		FromEmail: result.FromEmail,
		Date:      result.Date,
		Body:      result.Body,
	}

	if err := tx.Create(email).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	// 案件メールの場合、詳細情報を保存
	if result.MailCategory == "案件" {
		if err := r.saveProjectDetails(tx, result); err != nil {
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

// saveProjectDetails は案件メールの詳細情報を保存します
func (r *EmailStoreRepositoryImpl) saveProjectDetails(tx *gorm.DB, result *openaidomain.EmailAnalysisResult) error {
	// EmailProjectを保存
	emailProject := &domain.EmailProject{
		EmailID:             result.ID,
		MailCategory:        result.MailCategory,
		EndPeriod:           result.EndPeriod,
		WorkLocation:        result.WorkLocation,
		PriceFrom:           result.PriceFrom,
		PriceTo:             result.PriceTo,
		RemoteWorkCategory:  result.RemoteWorkCategory,
		RemoteWorkFrequency: result.RemoteWorkFrequency,
		// 一覧画面用のカンマ区切り文字列を作成
		TechnologiesText: r.createTechnologiesText(result),
		PositionsText:    "", // 今回は空文字列
		WorkTypesText:    "", // 今回は空文字列
	}

	if err := tx.Create(emailProject).Error; err != nil {
		return fmt.Errorf("EmailProject保存エラー: %w", err)
	}

	// EntryTimingを保存
	if err := r.saveEntryTimings(tx, emailProject.ID, result.StartPeriod); err != nil {
		return fmt.Errorf("EntryTiming保存エラー: %w", err)
	}

	// キーワード関連を保存
	if err := r.saveKeywords(tx, result); err != nil {
		return fmt.Errorf("キーワード保存エラー: %w", err)
	}

	return nil
}

// createTechnologiesText は技術要素のカンマ区切り文字列を作成します
func (r *EmailStoreRepositoryImpl) createTechnologiesText(result *openaidomain.EmailAnalysisResult) string {
	var technologies []string
	technologies = append(technologies, result.Languages...)
	technologies = append(technologies, result.Frameworks...)
	technologies = append(technologies, result.RequiredSkillsMust...)
	technologies = append(technologies, result.RequiredSkillsWant...)
	return strings.Join(technologies, ",")
}

// saveEntryTimings は入場時期を保存します
func (r *EmailStoreRepositoryImpl) saveEntryTimings(tx *gorm.DB, emailProjectID uint, startPeriods []string) error {
	for _, period := range startPeriods {
		entryTiming := &domain.EntryTiming{
			EmailProjectID: emailProjectID,
			Timing:         period,
		}
		if err := tx.Create(entryTiming).Error; err != nil {
			return fmt.Errorf("EntryTiming保存エラー: %w", err)
		}
	}
	return nil
}

// saveKeywords はキーワード関連のデータを保存します
func (r *EmailStoreRepositoryImpl) saveKeywords(tx *gorm.DB, result *openaidomain.EmailAnalysisResult) error {
	// 言語
	if err := r.saveKeywordsByType(tx, result.ID, result.Languages, "language"); err != nil {
		return err
	}

	// フレームワーク
	if err := r.saveKeywordsByType(tx, result.ID, result.Frameworks, "framework"); err != nil {
		return err
	}

	// 必須スキル
	if err := r.saveKeywordsByType(tx, result.ID, result.RequiredSkillsMust, "skill_must"); err != nil {
		return err
	}

	// 希望スキル
	if err := r.saveKeywordsByType(tx, result.ID, result.RequiredSkillsWant, "skill_want"); err != nil {
		return err
	}

	return nil
}

// saveKeywordsByType は指定されたタイプのキーワードを保存します
func (r *EmailStoreRepositoryImpl) saveKeywordsByType(tx *gorm.DB, emailID string, keywords []string, keywordType string) error {
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
			KeywordGroupID: keywordGroup.ID,
			Type:           keywordType,
		}

		if err := tx.Create(emailKeywordGroup).Error; err != nil {
			return fmt.Errorf("EmailKeywordGroup保存エラー: %w", err)
		}
	}
	return nil
}

// getOrCreateKeywordGroup はKeywordGroupを取得または作成します
func (r *EmailStoreRepositoryImpl) getOrCreateKeywordGroup(tx *gorm.DB, name string) (*domain.KeywordGroup, error) {
	var keywordGroup domain.KeywordGroup

	// 既存のKeywordGroupを検索
	err := tx.Where("name = ?", name).First(&keywordGroup).Error
	if err == nil {
		return &keywordGroup, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("KeywordGroup検索エラー: %w", err)
	}

	// 新規作成
	keywordGroup = domain.KeywordGroup{
		Name: name,
	}

	if err := tx.Create(&keywordGroup).Error; err != nil {
		return nil, fmt.Errorf("KeywordGroup作成エラー: %w", err)
	}

	return &keywordGroup, nil
}

// GetEmailByID はIDでメールを取得します
func (r *EmailStoreRepositoryImpl) GetEmailByID(ctx context.Context, id string) (*domain.Email, error) {
	var email domain.Email
	err := r.db.Where("id = ?", id).First(&email).Error
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
	err := r.db.Model(&domain.Email{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}
	return count > 0, nil
}
