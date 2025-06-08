// Package application はGメール機能群のアプリケーション層を提供します。
package application

import (
	cd "business/internal/common/domain"
	gi "business/internal/gmail/infrastructure"
	"context"
)

// GmailUseCase はGメール機能群のユースケースです
type GmailUseCase struct {
	r gi.GmailConnectInterface
}

// New は新しいメール機能群のユースケースを作成します
func New(r gi.GmailConnectInterface) *GmailUseCase {
	return &GmailUseCase{
		r: r,
	}
}

func (g *GmailUseCase) GetMessages(ctx context.Context, labelName string) ([]cd.BasicMessage, error) {
	return g.r.GetMessages(ctx, labelName)
}
