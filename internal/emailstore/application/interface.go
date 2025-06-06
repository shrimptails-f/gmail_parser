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

	// CheckGmailExists はGメールIDの存在チェックを行います
	CheckGmailIdExists(gmailId string) (bool, error)
}
