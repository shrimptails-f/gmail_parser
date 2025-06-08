// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	cd "business/internal/common/domain"
	r "business/internal/emailstore/infrastructure"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// EmailStoreUseCaseImpl はメール保存のユースケース実装です
type EmailStoreUseCaseImpl struct {
	r r.EmailStoreRepository
}

// NewEmailStoreUseCase はメール保存ユースケースを作成します
func NewEmailStoreUseCase(r r.EmailStoreRepository) *EmailStoreUseCaseImpl {
	return &EmailStoreUseCaseImpl{
		r: r,
	}
}

// SaveEmailAnalysisResult はメール分析結果を保存します
func (u *EmailStoreUseCaseImpl) SaveEmailAnalysisResult(result cd.Email) error {
	// リポジトリを使用してメールを保存
	if err := u.r.SaveEmail(result); err != nil {
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	return nil
}

// CheckGmailIdExists はメールIDの存在チェックを行います
func (u *EmailStoreUseCaseImpl) CheckGmailIdExists(emailId string) (bool, error) {
	if emailId == "" {
		return false, fmt.Errorf("メールIDが空です")
	}

	exists, err := u.r.EmailExists(emailId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("メール存在チェックエラー: %w", err)
	}

	return exists, nil
}
