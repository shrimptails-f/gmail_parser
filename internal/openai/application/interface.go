// Package application はAI機能のアプリケーション層を提供します。
// このファイルはAI機能のインターフェースを定義します。
package application

import (
	ed "business/internal/emailstore/domain"
	authdomain "business/internal/gmail/domain"
	"business/internal/openai/domain"
	"context"
)

// EmailAnalysisUseCase はメール分析のユースケースインターフェースです
type EmailAnalysisUseCase interface {
	AnalyzeEmailContent(ctx context.Context, message *authdomain.GmailMessage) (*domain.EmailAnalysisResult, error)
	DisplayEmailAnalysisResult(result *domain.EmailAnalysisResult) error
}

// TextAnalysisUseCase はテキスト字句解析のユースケースインターフェースです
type TextAnalysisUseCase interface {
	AnalyzeText(ctx context.Context, text string) (*domain.TextAnalysisResult, error)
	AnalyzeTextWithOptions(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error)
	AnalyzeEmailText(ctx context.Context, email authdomain.GmailMessage) ([]ed.AnalysisResult, error)
	AnalyzeEmailTextMultiple(ctx context.Context, emailText, messageID, subject string) ([]*domain.TextAnalysisResult, error)
	DisplayAnalysisResult(result *domain.TextAnalysisResult) error
}

// TextAnalysisService はテキスト字句解析サービスのインターフェースです
type TextAnalysisService interface {
	AnalyzeText(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error)
	AnalyzeTextMultiple(ctx context.Context, request *domain.TextAnalysisRequest) ([]*domain.TextAnalysisResult, error)
}
