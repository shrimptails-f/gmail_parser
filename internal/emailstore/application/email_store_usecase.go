// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	"business/internal/emailstore/domain"
	r "business/internal/emailstore/infrastructure"
	openaidomain "business/internal/openai/domain"
	"errors"
	"strings"

	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EmailStoreUseCaseImpl はメール保存のユースケース実装です
type EmailStoreUseCaseImpl struct {
	emailStoreRepository r.EmailStoreRepository
}

// NewEmailStoreUseCase はメール保存ユースケースを作成します
func NewEmailStoreUseCase(emailStoreRepository r.EmailStoreRepository) EmailStoreUseCase {
	return &EmailStoreUseCaseImpl{
		emailStoreRepository: emailStoreRepository,
	}
}

// SaveEmailAnalysisResult はメール分析結果を保存します
func (u *EmailStoreUseCaseImpl) SaveEmailAnalysisResult(ctx context.Context, result domain.AnalysisResult) error {
	// メールが既に存在するかチェック
	isGmailIdExist, err := u.CheckGmailIdExists(ctx, result.GmailID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("メール保存エラー: %w", err)
	}
	if isGmailIdExist {
		fmt.Printf("GメールID %s は既に処理済みです。字句解析をスキップします。\n", result.GmailID)
		return nil
	}

	// 言語のチェック
	notExistLanguages, err := u.filterNotExistingWords(result.Languages)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("言語の存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups := u.buildKeywordEntities(notExistLanguages, "LANGUAGE")

	// フレームワークのチェック
	notExistFrameworks, err := u.filterNotExistingWords(result.Frameworks)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("フレームワークの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistFrameworks, "FRAMEWORK")...)

	// 必須スキルのチェック
	notExistRequiredSkillsMust, err := u.filterNotExistingWords(result.RequiredSkillsMust)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("必須スキルの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistRequiredSkillsMust, "MUST")...)

	// 尚可のチェック
	notExistRequiredSkillsWant, err := u.filterNotExistingWords(result.RequiredSkillsWant)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("尚可スキルの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistRequiredSkillsWant, "WANT")...)

	// ポジションチェック
	notExistPositionWords, err := u.filterNotExistingPositionWords(result.Positions)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("ポジションの存在確認でエラーが発生しました。")
		return err
	}
	emailPositionGroups := u.buildPositionEntities(notExistPositionWords)

	// 業種チェック
	notExistWorkTypeWords, err := u.filterNotExistingWorkTypeWords(result.WorkTypes)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("ポジションの存在確認でエラーが発生しました。")
		return err
	}
	emailWorkTypeGroups := u.buildWorkTypeEntities(notExistWorkTypeWords)

	// 引数の情報を保存用の構造体に詰め替える。
	email := convertToProjectAnalysisResult(result, emailKeywordGroups, emailPositionGroups, emailWorkTypeGroups)

	// リポジトリを使用してメールを保存
	if err := u.emailStoreRepository.SaveEmail(email); err != nil {
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	return nil
}

// convertToProjectAnalysisResult はTextAnalysisResultを案件分析結果に変換します
func convertToProjectAnalysisResult(result domain.AnalysisResult, emailKeywordGroup []r.EmailKeywordGroup, emailPositionGroup []r.EmailPositionGroup, emailWorkTypeGroups []r.EmailWorkTypeGroup) r.Email {
	sep := ","
	emailProject := r.EmailProject{
		ProjectTitle:    stringPtr(result.ProjectName),
		EntryTiming:     stringPtr(strings.Join(result.StartPeriod, sep)),
		EndTiming:       stringPtr(result.EndPeriod),
		WorkLocation:    stringPtr(result.WorkLocation),
		PriceFrom:       result.PriceFrom,
		PriceTo:         result.PriceTo,
		Languages:       stringPtr(strings.Join(result.Languages, sep)),
		Frameworks:      stringPtr(strings.Join(result.Frameworks, sep)),
		Positions:       stringPtr(strings.Join(result.Positions, sep)),
		WorkTypes:       stringPtr(strings.Join(result.WorkTypes, sep)),
		MustSkills:      stringPtr(strings.Join(result.RequiredSkillsMust, sep)),
		WantSkills:      stringPtr(strings.Join(result.RequiredSkillsWant, sep)),
		RemoteType:      result.RemoteWorkCategory,
		RemoteFrequency: result.RemoteWorkFrequency,
	}

	email := r.Email{
		GmailID:      result.GmailID,
		Subject:      result.Subject,
		SenderName:   result.SenderName(),
		SenderEmail:  result.SenderEmail(),
		ReceivedDate: time.Time{},
		Body:         stringPtr(result.Body),
		Category:     result.Category,

		EmailProject: &emailProject,
	}

	return email
}

// SaveEmailAnalysisMultipleResult は複数案件対応のメール分析結果を保存します
func (u *EmailStoreUseCaseImpl) SaveEmailAnalysisMultipleResult(ctx context.Context, result *openaidomain.EmailAnalysisMultipleResult) error {
	// 入力値チェック
	if result == nil {
		return fmt.Errorf("分析結果がnilです")
	}

	// 結果の妥当性チェック
	if err := result.IsValid(); err != nil {
		return fmt.Errorf("分析結果妥当性チェックエラー: %w", err)
	}

	// リポジトリを使用してメールを保存
	if err := u.emailStoreRepository.SaveEmailMultiple(ctx, result); err != nil {
		return fmt.Errorf("複数案件メール保存エラー: %w", err)
	}

	return nil
}

// CheckGmailIdExists はメールIDの存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckGmailIdExists(ctx context.Context, emailId string) (bool, error) {
	if emailId == "" {
		return false, fmt.Errorf("メールIDが空です")
	}

	exists, err := u.emailStoreRepository.EmailExists(ctx, emailId)
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}

	return exists, nil
}

// filterNotExistingWords はキーワード取得または存在確認を行い、作成が必要なキーワード一覧を返します
func (u *EmailStoreUseCaseImpl) filterNotExistingWords(words []string) ([]string, error) {
	keywords, err := u.emailStoreRepository.GetKeywords(words)
	if err != nil {
		return nil, err
	}

	existing := make(map[string]bool)
	for _, keyword := range keywords {
		existing[keyword.Word] = true
	}
	// 取得した単語数と引数の単語数が一致したら確認処理を終了する。
	if len(words) == len(existing) {
		return nil, nil
	}

	var notExistWordsTmp []string
	for _, word := range words {
		if !existing[word] {
			notExistWordsTmp = append(notExistWordsTmp, word)
		}
	}

	keywordGroups, err := u.emailStoreRepository.GetkeywordGroups(notExistWordsTmp)
	if err != nil {
		return nil, err
	}

	existing = make(map[string]bool)
	for _, keywordGroup := range keywordGroups {
		existing[keywordGroup.Name] = true
	}

	var notExistWords []string
	for _, word := range words {
		if !existing[word] {
			notExistWords = append(notExistWords, word)
		}
	}

	return notExistWords, nil
}

func (u *EmailStoreUseCaseImpl) buildKeywordEntities(words []string, keywordType string) []r.EmailKeywordGroup {
	now := time.Now()
	var emailKeywordGroups []r.EmailKeywordGroup

	for _, word := range words {
		emailKeywordGroup := r.EmailKeywordGroup{
			CreatedAt: now,
			Type:      strings.ToUpper(keywordType), // 例: "must" → "MUST"
			KeywordGroup: r.KeywordGroup{
				Name:      word,
				CreatedAt: now,
				UpdatedAt: now,
				WordLinks: []r.KeywordGroupWordLink{
					{
						KeyWord: r.KeyWord{
							Word:      word,
							CreatedAt: now,
							UpdatedAt: now,
						},
						CreatedAt: now,
					},
				},
			},
		}

		emailKeywordGroups = append(emailKeywordGroups, emailKeywordGroup)
	}

	return emailKeywordGroups
}

func (u *EmailStoreUseCaseImpl) filterNotExistingPositionWords(words []string) ([]string, error) {
	positionWords, err := u.emailStoreRepository.GetPositionWords(words)
	if err != nil {
		return nil, err
	}

	existing := make(map[string]bool)
	for _, pw := range positionWords {
		existing[pw.Word] = true
	}
	// 取得した単語数と引数の単語数が一致したら確認処理を終了する。
	if len(words) == len(existing) {
		return nil, nil
	}

	var notExistWordsTmp []string
	for _, word := range words {
		if !existing[word] {
			notExistWordsTmp = append(notExistWordsTmp, word)
		}
	}

	positionGroups, err := u.emailStoreRepository.GetPositionGroups(notExistWordsTmp)
	if err != nil {
		return nil, err
	}

	existing = make(map[string]bool)
	for _, group := range positionGroups {
		existing[group.Name] = true
	}

	var notExistWords []string
	for _, word := range words {
		if !existing[word] {
			notExistWords = append(notExistWords, word)
		}
	}

	return notExistWords, nil
}

func (u *EmailStoreUseCaseImpl) filterNotExistingWorkTypeWords(words []string) ([]string, error) {
	workTypeWords, err := u.emailStoreRepository.GetWorkTypeWords(words)
	if err != nil {
		return nil, err
	}

	existing := make(map[string]bool)
	for _, w := range workTypeWords {
		existing[w.Word] = true
	}

	if len(words) == len(existing) {
		return nil, nil
	}

	var notExistWordsTmp []string
	for _, word := range words {
		if !existing[word] {
			notExistWordsTmp = append(notExistWordsTmp, word)
		}
	}

	workTypeGroups, err := u.emailStoreRepository.GetWorkTypeGroups(notExistWordsTmp)
	if err != nil {
		return nil, err
	}

	existing = make(map[string]bool)
	for _, group := range workTypeGroups {
		existing[group.Name] = true
	}

	var notExistWords []string
	for _, word := range words {
		if !existing[word] {
			notExistWords = append(notExistWords, word)
		}
	}

	return notExistWords, nil
}

func (u *EmailStoreUseCaseImpl) buildPositionEntities(words []string) []r.EmailPositionGroup {
	var emailPositionGroups []r.EmailPositionGroup

	for _, word := range words {
		emailPositionGroup := r.EmailPositionGroup{
			PositionGroup: r.PositionGroup{
				Name: word,
				Words: []r.PositionWord{
					{
						Word: word,
					},
				},
			},
		}
		emailPositionGroups = append(emailPositionGroups, emailPositionGroup)
	}

	return emailPositionGroups
}

func (u *EmailStoreUseCaseImpl) buildWorkTypeEntities(words []string) []r.EmailWorkTypeGroup {
	var emailWorkTypeGroups []r.EmailWorkTypeGroup

	for _, word := range words {
		emailWorkTypeGroup := r.EmailWorkTypeGroup{
			WorkTypeGroup: r.WorkTypeGroup{
				Name: word,
				Words: []r.WorkTypeWord{
					{
						Word: word,
					},
				},
			},
		}
		emailWorkTypeGroups = append(emailWorkTypeGroups, emailWorkTypeGroup)
	}

	return emailWorkTypeGroups
}

// // CheckKeywordExists はキーワードの存在チェックを行います
// func (u *EmailStoreUseCaseImpl) CheckKeywordExists(ctx context.Context, word string) (bool, error) {
// 	if word == "" {
// 		return false, fmt.Errorf("キーワードが空です")
// 	}

// 	exists, err := u.emailStoreRepository.GetKeywords(words)
// 	if err != nil {
// 		return false, fmt.Errorf("キーワード存在チェックエラー: %w", err)
// 	}

// 	return exists, nil
// }

func stringPtr(s string) *string { return &s }
