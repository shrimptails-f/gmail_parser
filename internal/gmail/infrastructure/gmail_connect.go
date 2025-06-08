// Package infrastructure はGメールとの疎通部分を実装します。
package infrastructure

import (
	cd "business/internal/common/domain"
	gc "business/tools/gmail"
	"context"
)

// GmailConnect Gメール接続処理を持つ構造体です。
type GmailConnect struct {
	gc gc.ClientInterface
}

func New(gc gc.ClientInterface) *GmailConnect {
	return &GmailConnect{
		gc: gc,
	}
}

func (g *GmailConnect) GetMessages(ctx context.Context, labelName string) ([]cd.BasicMessage, error) {
	return g.gc.GetMessagesByLabelName(ctx, labelName)
}
