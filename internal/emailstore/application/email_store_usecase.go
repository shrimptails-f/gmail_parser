// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	openaidomain "business/internal/openai/domain"
	"context"
	"fmt"
)

// EmailStoreUseCaseImpl はメール保存のユースケース実装です
type EmailStoreUseCaseImpl struct {
	emailStoreRepository EmailStoreRepository
}

// NewEmailStoreUseCase はメール保存ユースケースを作成します
func NewEmailStoreUseCase(emailStoreRepository EmailStoreRepository) EmailStoreUseCase {
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
