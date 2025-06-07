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

	// GetkeywordGroups は引数からまとめて取得します。
	GetkeywordGroups(name []string) ([]KeywordGroup, error)

	// GetKeywords は引数からまとめて取得します。
	GetKeywords(words []string) ([]KeyWord, error)

	// GetPositionGroups は引数からまとめて取得します。
	GetPositionGroups(name []string) ([]PositionGroup, error)

	// GetPositionWords はポジションが既に存在するかチェックします
	GetPositionWords(words []string) ([]PositionWord, error)

	// WorkTypeExists は業務種別が既に存在するかチェックします
	GetWorkTypeWords(words []string) ([]WorkTypeWord, error)

	// GetWorkTypeGroups は業務種別が既に存在するかチェックします
	GetWorkTypeGroups(words []string) ([]WorkTypeGroup, error)
}
