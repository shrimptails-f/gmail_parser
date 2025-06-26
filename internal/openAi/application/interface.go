// Package application はAI機能のアプリケーション層を提供します。
// このファイルはAI機能のインターフェースを定義します。
package application

import (
	cd "business/internal/common/domain"
	"context"
)

// EmailAnalysisUseCase はメール分析のユースケースインターフェースです
type EmailAnalysisUseCase interface {
	AnalyzeEmailContent(ctx context.Context, emails []cd.BasicMessage) ([]cd.Email, error)
}
