package presentation

import (
	cd "business/internal/common/domain"
	ea "business/internal/emailstore/application"
	ga "business/internal/gmail/application"
	aiapp "business/internal/openAi/application"
	"context"
	"fmt"
)

// AnalyzeEmailController はメール保存のユースケース実装です
type AnalyzeEmailController struct {
	ea    ea.UseCaseInterface
	ga    ga.UseCaseInterface
	aiapp aiapp.UseCaseInterface
}

// New はメール保存ユースケースを作成します
func New(
	ea ea.UseCaseInterface,
	ga ga.UseCaseInterface,
	aiapp aiapp.UseCaseInterface,
) *AnalyzeEmailController {
	return &AnalyzeEmailController{
		ea:    ea,
		ga:    ga,
		aiapp: aiapp,
	}
}

func (n *AnalyzeEmailController) SaveEmailAnalysisResult(ctx context.Context) error {
	messages, err := n.ga.GetMessages(ctx, "営業/案件", -1)
	if err != nil {
		fmt.Printf("gメール取得処理失敗: %v \n", err)
		return err
	}
	if len(messages) == 0 {
		fmt.Printf("gメールの取得結果が0件だったため処理を終了しました。\n")
		return nil
	}

	fmt.Printf("メール分析を行います。 \n")
	var analysisResults []cd.Email
	analysisResults, err = n.aiapp.AnalyzeEmailContent(ctx, messages)
	if err != nil {
		fmt.Printf("メール分析エラー: %v \n", err)
		return err
	}

	fmt.Printf("DBへの保存処理を開始します。")
	for _, email := range analysisResults {
		err = n.ea.SaveEmailAnalysisResult(email)
		if err != nil {
			fmt.Printf("メール保存エラー: %v \n", err)
			return err
		}
	}
	fmt.Printf("DBへの保存処理が完了しました。 \n")
	return nil
}
