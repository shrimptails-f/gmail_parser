// Package infrastructure はAI機能のアプリケーション層を提供します。
// このファイルはAI機能のインターフェースを定義します。
package infrastructure

import (
	"business/internal/openai/domain"
	"context"
)

// EmailAnalysisService はメール分析サービスのインターフェースです
type EmailAnalysisService interface {
	AnalyzeEmail(ctx context.Context, request *domain.EmailAnalysisRequest) (*domain.EmailAnalysisResult, error)
}

// PromptService はプロンプト管理サービスのインターフェースです
type PromptService interface {
	LoadPrompt(filename string) (string, error)
	SavePrompt(filename, content string) error
	ListPrompts() ([]string, error)
}
