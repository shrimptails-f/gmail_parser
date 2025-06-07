// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	cd "business/internal/common/domain"
)

// EmailStoreRepository はメール保存のリポジトリインターフェースです
type EmailStoreRepository interface {
	// SaveEmail はメール分析結果をデータベースに保存します
	SaveEmail(result cd.Email) error

	// GetEmailByGmailId はIDでメールを取得します
	GetEmailByGmailId(gmail_id string) (Email, error)

	// EmailExists はメールが既に存在するかチェックします
	EmailExists(id string) (bool, error)
}
