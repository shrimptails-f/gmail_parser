// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するユースケースを実装します。
package application

import (
	openaidomain "business/internal/openai/domain"
	"context"
)

// EmailStoreUseCase はメール保存のユースケースインターフェースです
type EmailStoreUseCase interface {
	// SaveEmailAnalysisResult はメール分析結果を保存します
	SaveEmailAnalysisResult(ctx context.Context, result *openaidomain.EmailAnalysisResult) error

	// SaveEmailAnalysisMultipleResult は複数案件対応のメール分析結果を保存します
	SaveEmailAnalysisMultipleResult(ctx context.Context, result *openaidomain.EmailAnalysisMultipleResult) error

	// CheckGmailExists はGメールIDの存在チェックを行います
	CheckGmailIdExists(ctx context.Context, gmailId string) (bool, error)

	// CheckKeywordExists はキーワードの存在チェックを行います
	CheckKeywordExists(ctx context.Context, word string) (bool, error)

	// CheckPositionExists はポジションの存在チェックを行います
	CheckPositionExists(ctx context.Context, word string) (bool, error)

	// CheckWorkTypeExists は業務種別の存在チェックを行います
	CheckWorkTypeExists(ctx context.Context, word string) (bool, error)
}
