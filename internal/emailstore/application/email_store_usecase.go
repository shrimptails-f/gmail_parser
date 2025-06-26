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

// UseCase はメール保存のユースケースの具象です
type UseCase struct {
	r r.RepositoryInterface
}

// New はメール保存ユースケースを作成します
func New(r r.RepositoryInterface) *UseCase {
	return &UseCase{
		r: r,
	}
}

// SaveEmailAnalysisResult はメール分析結果を保存します
func (u *UseCase) SaveEmailAnalysisResult(result cd.Email) error {
	// リポジトリを使用してメールを保存
	if err := u.r.SaveEmail(result); err != nil {
		return fmt.Errorf("メール保存エラー: %w", err)
	}

	return nil
}

// GetEmailByGmailIds はメールIDリストを返却します
func (u *UseCase) GetEmailByGmailIds(emailIdList []string) ([]string, error) {
	if len(emailIdList) == 0 {
		return []string{}, nil
	}

	exists, err := u.r.GetEmailByGmailIds(emailIdList)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return []string{}, fmt.Errorf("メール存在チェックエラー: %w", err)
	}

	return exists, nil
}
