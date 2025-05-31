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
	DisplayEmailAnalysisResult(result *domain.EmailAnalysisResult) error
}

// EmailAnalysisService はメール分析サービスのインターフェースです
type EmailAnalysisService interface {
	AnalyzeEmail(ctx context.Context, request *domain.EmailAnalysisRequest) (*domain.EmailAnalysisResult, error)
}

// TextAnalysisUseCase はテキスト字句解析のユースケースインターフェースです
type TextAnalysisUseCase interface {
	AnalyzeText(ctx context.Context, text string) (*domain.TextAnalysisResult, error)
	AnalyzeTextWithOptions(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error)
	DisplayAnalysisResult(result *domain.TextAnalysisResult) error
}

// TextAnalysisService はテキスト字句解析サービスのインターフェースです
type TextAnalysisService interface {
	AnalyzeText(ctx context.Context, request *domain.TextAnalysisRequest) (*domain.TextAnalysisResult, error)
}

// PromptService はプロンプト管理サービスのインターフェースです
type PromptService interface {
	LoadPrompt(filename string) (string, error)
	SavePrompt(filename, content string) error
	ListPrompts() ([]string, error)
}
