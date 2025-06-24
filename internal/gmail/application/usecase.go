// Package application はGメール機能群のアプリケーション層を提供します。
package application

import (
	cd "business/internal/common/domain"
	ea "business/internal/emailstore/application"
	gi "business/internal/gmail/infrastructure"
	"context"
	"fmt"
	"sync"

	"github.com/samber/lo"
)

// GmailUseCase はGメール機能群のユースケースです
type GmailUseCase struct {
	r  gi.GmailConnectInterface
	ea ea.EmailStoreUseCase
}

// New は新しいメール機能群のユースケースを作成します
func New(r gi.GmailConnectInterface, ea ea.EmailStoreUseCase) *GmailUseCase {
	return &GmailUseCase{
		r:  r,
		ea: ea,
	}
}

func (g *GmailUseCase) GetMessages(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error) {
	ids, err := g.r.GetMessageIds(ctx, labelName, sinceDaysAgo)
	if err != nil {
		return nil, err
	}
	fmt.Printf("取得したメッセージ数: %d\n\n", len(ids))

	getIds, err := g.ea.GetEmailByGmailIds(ids)
	if err != nil {
		return nil, fmt.Errorf("GetMessages: %v", err)
	}

	for _, id := range lo.Intersect(ids, getIds) {
		fmt.Printf("GメールID: %v は登録済みのため解析をスキップしました。 \n", id)
	}

	// getIdsに存在しないIDを取得 つまりDBに登録する必要のあるメールということ。
	notExistIds, _ := lo.Difference(ids, getIds)
	var checkExistsWg sync.WaitGroup
	existMessagesChan := make(chan cd.BasicMessage, len(notExistIds))
	for _, id := range notExistIds {
		checkExistsWg.Add(1)

		go func(messageId string) {
			defer checkExistsWg.Done()

			email, err := g.r.GetGmailDetail(messageId)
			if err != nil {
				fmt.Printf("Gメール詳細取得時にエラーが発生しました。: %v\n", err)
				return
			}

			existMessagesChan <- email
		}(id)
	}
	checkExistsWg.Wait()
	close(existMessagesChan)

	var existMessages []cd.BasicMessage
	for msg := range existMessagesChan {
		existMessages = append(existMessages, msg)
	}

	return existMessages, nil
}
