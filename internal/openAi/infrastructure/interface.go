// Package infrastructure はAI機能のアプリケーション層を提供します。
// このファイルはAI機能のインターフェースを定義します。
package infrastructure

import (
	cd "business/internal/common/domain"
	"context"
)

// ConnectInterface はOpenAI APIのインターフェースです。
type ConnectInterface interface {
	AnalyzeEmailBody(ctx context.Context, prompt string) ([]cd.AnalysisResult, error)
}
