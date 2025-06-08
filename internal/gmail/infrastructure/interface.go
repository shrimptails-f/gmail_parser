// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルは認証機能で使用するインターフェースを定義します。
package infrastructure

import (
	cd "business/internal/common/domain"
	"context"
)

// GmailConnectInterface はGメール取得のインターフェースです。
type GmailConnectInterface interface {
	// GetMessages はラベルからメールを取得します。
	GetMessages(ctx context.Context, labelName string, sinceDaysAgo int) ([]cd.BasicMessage, error)
}
