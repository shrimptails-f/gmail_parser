// Package application はメール分析のアプリケーション層を提供します。
// このファイルはメール分析に関するユースケースを実装します。
package application

import (
	cd "business/internal/common/domain"
	r "business/internal/openAi/infrastructure"
	"business/tools/oswrapper"
	"context"
	"fmt"
	"sync"
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
func (u *UseCase) AnalyzeEmailContent(ctx context.Context, emails []cd.BasicMessage) ([]cd.Email, error) {
	// TODO あとでENVに追加する。
	prompt, err := u.os.ReadFile("/data/prompts/text_analysis_prompt.txt")
	if err != nil {
		return nil, err
	}

	var AnalyzeEmailWg sync.WaitGroup
	analyzeEmailChan := make(chan cd.Email)
	for _, email := range emails {
		AnalyzeEmailWg.Add(1)

		// メール本文の分析を実行
		go func(email cd.BasicMessage) {
			defer AnalyzeEmailWg.Done()

			analysisResults, err := u.r.AnalyzeEmailBody(ctx, string(prompt)+"\n\n"+email.Body)

			if err != nil {
				fmt.Printf("解析時にエラーが発生しました。 GメールID: %s %v \n", email.ID, err)
				return
			}
			if len(analysisResults) == 0 {
				fmt.Printf("GメールID: %v の解析結果が0件でした。 メールを確認してください。\n", email.ID)
				return
			}

			// 解析結果を保存形式へ詰め替える。
			results := convertToStructs(email, analysisResults)
			for _, result := range results {
				analyzeEmailChan <- result
			}
		}(email)
	}
	go func() { AnalyzeEmailWg.Wait(); close(analyzeEmailChan) }()

	var analysisEmail []cd.Email
	for email := range analyzeEmailChan {
		analysisEmail = append(analysisEmail, email)
	}

	return analysisEmail, nil
}

// convertToStructs は引数を結合して保存する形式へ詰め替えます。
func convertToStructs(message cd.BasicMessage, analysisResults []cd.AnalysisResult) []cd.Email {
	var results []cd.Email

	for _, analysisResult := range analysisResults {
		result := cd.Email{
			GmailID:             message.ID,
			ReceivedDate:        message.Date,
			Summary:             analysisResult.ProjectTitle,
			Subject:             message.Subject,
			From:                message.From,
			FromEmail:           message.ExtractEmailAddress(),
			Body:                message.Body,
			Category:            analysisResult.MailCategory,
			ProjectName:         analysisResult.ProjectTitle,
			StartPeriod:         analysisResult.StartPeriod,
			EndPeriod:           analysisResult.EndPeriod,
			WorkLocation:        analysisResult.WorkLocation,
			PriceFrom:           analysisResult.PriceFrom,
			PriceTo:             analysisResult.PriceTo,
			Languages:           analysisResult.Languages,
			Frameworks:          analysisResult.Frameworks,
			Positions:           analysisResult.Positions,
			WorkTypes:           analysisResult.WorkTypes,
			RequiredSkillsMust:  analysisResult.RequiredSkillsMust,
			RequiredSkillsWant:  analysisResult.RequiredSkillsWant,
			RemoteWorkCategory:  analysisResult.RemoteWorkCategory,
			RemoteWorkFrequency: analysisResult.RemoteWorkFrequency,
		}
		results = append(results, result)
	}

	return results
}
