// Package infrastructure は認証機能のインフラストラクチャ層を提供します。
// このファイルは認証機能で使用するインターフェースを定義します。
package infrastructure

import (
	cd "business/internal/common/domain"
	"context"
)

// GmailConnectInterface はGメール取得のインターフェースです。
type GmailConnectInterface interface {
	// GetMessageIds はラベルからメールを取得します。
	GetMessageIds(ctx context.Context, labelName string, sinceDaysAgo int) ([]string, error)
	// GetGmailDetail はIDからGメールを取得します。
	GetGmailDetail(id string) (cd.BasicMessage, error)
}
