// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	"business/internal/emailstore/domain"
	r "business/internal/emailstore/infrastructure"
	"errors"
	"fmt"
	"strings"

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
func (u *EmailStoreUseCaseImpl) SaveEmailAnalysisResult(result domain.AnalysisResult) error {
	// メールが既に存在するかチェック
	isGmailIdExist, err := u.CheckGmailIdExists(result.GmailID)
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
	emailKeywordGroups := u.buildKeywordEntities(notExistLanguages, "language")

	// フレームワークのチェック
	notExistFrameworks, err := u.filterNotExistingWords(result.Frameworks)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("フレームワークの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistFrameworks, "framework")...)

	// 必須スキルのチェック
	notExistRequiredSkillsMust, err := u.filterNotExistingWords(result.RequiredSkillsMust)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("必須スキルの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistRequiredSkillsMust, "must")...)

	// 尚可のチェック
	notExistRequiredSkillsWant, err := u.filterNotExistingWords(result.RequiredSkillsWant)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("尚可スキルの存在確認でエラーが発生しました。")
		return err
	}
	emailKeywordGroups = append(emailKeywordGroups, u.buildKeywordEntities(notExistRequiredSkillsWant, "want")...)

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

// convertToProjectAnalysisResult はDB保存構造体へ詰め替えます
func convertToProjectAnalysisResult(result domain.AnalysisResult, emailKeywordGroups []r.EmailKeywordGroup, emailPositionGroups []r.EmailPositionGroup, emailWorkTypeGroups []r.EmailWorkTypeGroup) r.Email {
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

	fmt.Printf("%v \n", result.ReceivedDate)
	email := r.Email{
		GmailID:      result.GmailID,
		Subject:      result.Subject,
		SenderName:   result.SenderName(),
		SenderEmail:  result.SenderEmail(),
		ReceivedDate: result.ReceivedDate,
		Body:         stringPtr(result.Body),
		Category:     result.Category,

		EmailProject:        &emailProject,
		EmailKeywordGroups:  emailKeywordGroups,
		EmailPositionGroups: emailPositionGroups,
		EmailWorkTypeGroups: emailWorkTypeGroups,
	}

	return email
}

// CheckGmailIdExists はメールIDの存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckGmailIdExists(emailId string) (bool, error) {
	if emailId == "" {
		return false, fmt.Errorf("メールIDが空です")
	}

	exists, err := u.emailStoreRepository.EmailExists(emailId)
	if err != nil {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}

	return exists, nil
}

// filterNotExistingWords はキーワード取得または存在確認を行い、作成が必要なキーワード一覧を返します
func (u *EmailStoreUseCaseImpl) filterNotExistingWords(words []string) ([]string, error) {
	// TODO: typeをクエリに含める?
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
	var emailKeywordGroups []r.EmailKeywordGroup

	for _, word := range words {
		emailKeywordGroup := r.EmailKeywordGroup{
			Type: keywordType,
			KeywordGroup: r.KeywordGroup{
				Type: keywordType,
				Name: word,
				WordLinks: []r.KeywordGroupWordLink{
					{
						KeyWord: r.KeyWord{
							Word: word,
						},
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

func stringPtr(s string) *string { return &s }
