// Package application はGメール機能群のアプリケーション層を提供します。
package application

import (
	cd "business/internal/common/domain"
	ea "business/internal/emailstore/application"
	gi "business/internal/gmail/infrastructure"
	"context"
	"fmt"
	"sync"
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
	messages, err := g.r.GetMessages(ctx, labelName, sinceDaysAgo)
	fmt.Printf("取得したメッセージ数: %d\n\n", len(messages))

	if err != nil {
		return nil, err
	}

	var checkExistsWg sync.WaitGroup
	existMessagesChan := make(chan cd.BasicMessage, len(messages))
	for i, message := range messages {
		// 今回データが入れ替わるかもためしたいので、クロージャー変数は定義しない。
		checkExistsWg.Add(1)

		go func() {
			defer checkExistsWg.Done()

			var exists bool
			exists, err = g.ea.CheckGmailIdExists(message.ID)
			if err != nil {
				fmt.Printf("メール存在確認エラー: %v\n", err)
				return
			}
			if exists {
				fmt.Printf("%v通目 メールID %s は既に処理済みです。字句解析をスキップします。\n", i, message.ID)
				return
			}

			existMessagesChan <- message
		}()
	}
	checkExistsWg.Wait()
	close(existMessagesChan)

	count := 0
	var existMessages []cd.BasicMessage
	for msg := range existMessagesChan {
		count++
		if count < 5 {
			existMessages = append(existMessages, msg)
		} else {
			break // TODO: 消す
		}
	}

	return existMessages, nil
}
