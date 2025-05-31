// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	r "business/internal/emailstore/infrastructure"
	openaidomain "business/internal/openai/domain"
	"context"
	"fmt"
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
func (u *EmailStoreUseCaseImpl) SaveEmailAnalysisResult(ctx context.Context, result *openaidomain.EmailAnalysisResult) error {
	// 入力値チェック
	if result == nil {
		return fmt.Errorf("分析結果がnilです")
	}

	// リポジトリを使用してメールを保存
	if err := u.emailStoreRepository.SaveEmail(ctx, result); err != nil {
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	return nil
}

// CheckKeywordExists はキーワードの存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckKeywordExists(ctx context.Context, word string) (bool, error) {
	if word == "" {
		return false, fmt.Errorf("キーワードが空です")
	}

	exists, err := u.emailStoreRepository.KeywordExists(ctx, word)
	if err != nil {
		return false, fmt.Errorf("キーワード存在チェックエラー: %w", err)
	}

	return exists, nil
}

// CheckPositionExists はポジションの存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckPositionExists(ctx context.Context, word string) (bool, error) {
	if word == "" {
		return false, fmt.Errorf("ポジションが空です")
	}

	exists, err := u.emailStoreRepository.PositionExists(ctx, word)
	if err != nil {
		return false, fmt.Errorf("ポジション存在チェックエラー: %w", err)
	}

	return exists, nil
}

// CheckWorkTypeExists は業務種別の存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckWorkTypeExists(ctx context.Context, word string) (bool, error) {
	if word == "" {
		return false, fmt.Errorf("業務種別が空です")
	}

	exists, err := u.emailStoreRepository.WorkTypeExists(ctx, word)
	if err != nil {
		return false, fmt.Errorf("業務種別存在チェックエラー: %w", err)
	}

	return exists, nil
}
