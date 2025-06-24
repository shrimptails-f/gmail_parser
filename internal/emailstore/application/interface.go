// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	cd "business/internal/common/domain"
)

// EmailStoreUseCase はメール保存のユースケースインターフェースです
type EmailStoreUseCase interface {
	// SaveEmailAnalysisResult はメール分析結果を保存します
	SaveEmailAnalysisResult(result cd.Email) error

	// GetEmailByGmailIds はGメールIDリストを返却します。
	GetEmailByGmailIds(gmailId []string) ([]string, error)
}
