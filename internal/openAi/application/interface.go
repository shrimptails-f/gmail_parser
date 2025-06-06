// Package application はAI機能のアプリケーション層を提供します。
// このファイルはAI機能のインターフェースを定義します。
package application

import (
	authdomain "business/internal/gmail/domain"
	"business/internal/openai/domain"
	"context"
)

// EmailAnalysisUseCase はメール分析のユースケースインターフェースです
type EmailAnalysisUseCase interface {
	AnalyzeEmailContent(ctx context.Context, message *authdomain.GmailMessage) (*domain.EmailAnalysisResult, error)
}
