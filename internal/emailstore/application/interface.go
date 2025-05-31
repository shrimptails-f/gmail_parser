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
}
