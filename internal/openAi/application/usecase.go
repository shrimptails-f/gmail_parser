// Package application はメール分析のアプリケーション層を提供します。
// このファイルはメール分析に関するユースケースを実装します。
package application

import (
	cd "business/internal/common/domain"
	r "business/internal/openAi/infrastructure"
	"business/tools/oswrapper"
	"context"
)

// UseCase はメール分析のユースケース実装です
type UseCase struct {
	r  r.AnalyzerInterFace
	os oswrapper.OsWapperInterface
}

// NewEmailAnalysisUseCase はメール分析ユースケースを作成します
func NewUseCase(r r.AnalyzerInterFace, os oswrapper.OsWapperInterface) *UseCase {
	return &UseCase{
		r:  r,
		os: os,
	}
}

// AnalyzeEmailContent はメール内容を分析します
func (u *UseCase) AnalyzeEmailContent(ctx context.Context, mailBody string) ([]cd.AnalysisResult, error) {
	// TODO あとでENVに追加する。
	prompt, err := u.os.ReadFile("/data/prompts/text_analysis_prompt.txt")
	if err != nil {
		return nil, err
	}

	return u.r.AnalyzeEmailBody(ctx, string(prompt)+"\n\n"+mailBody)
}
