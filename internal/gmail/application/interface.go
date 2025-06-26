// Package application はGメール機能群のアプリケーション層を提供します。
package application

import (
	cd "business/internal/common/domain"
	"context"
)

// UseCaseInterface はGmailのユースケースインターフェースです
type UseCaseInterface interface {
	GetMessages(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error)
}
