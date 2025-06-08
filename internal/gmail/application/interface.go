// Package application はGメール機能群のアプリケーション層を提供します。
package application

import (
	cd "business/internal/common/domain"
	"context"
)

// GmailUseCaseInterface はGmailのユースケースインターフェースです
type GmailUseCaseInterface interface {
	GetMessages(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error)
}
