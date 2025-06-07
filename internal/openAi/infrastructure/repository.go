// Package infrastructure はAI機能のインフラストラクチャ層を提供します。
// このファイルはOpenAI APIとの通信を行うサービスを実装します。
package infrastructure

import (
	cd "business/internal/common/domain"
	oa "business/tools/openai"
	"context"
)

// UseCase はメール分析のユースケース実装です
type Analyzer struct {
	oa oa.ClientInterFace
}

// NewEmailAnalysisUseCase はメール分析ユースケースを作成します
func NewAnalyzer(oa oa.ClientInterFace) *Analyzer {
	return &Analyzer{
		oa: oa,
	}
}

// AnalyzeEmailBody はメール内容を分析します
func (u *Analyzer) AnalyzeEmailBody(ctx context.Context, prompt string) ([]cd.AnalysisResult, error) {
	return u.oa.Chat(ctx, prompt)
}
