// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	cd "business/internal/common/domain"
)

// RepositoryInterface はメール保存のリポジトリインターフェースです
type RepositoryInterface interface {
	// SaveEmail はメール分析結果をデータベースに保存します
	SaveEmail(result cd.Email) error

	// GetEmailByGmailIds はIDでメールを取得します
	GetEmailByGmailIds(gmail_ids []string) ([]string, error)
}
