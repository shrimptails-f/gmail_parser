// Package infrastructure はメール保存機能のインフラストラクチャ層を提供します。
// このファイルはメール保存に関するリポジトリの実装を提供します。
package infrastructure

import (
	"business/internal/emailstore/domain"
	openaidomain "business/internal/openai/domain"
	"context"
)

// EmailStoreRepository はメール保存のリポジトリインターフェースです
type EmailStoreRepository interface {
	// SaveEmail はメール分析結果をデータベースに保存します
	SaveEmail(ctx context.Context, result *openaidomain.EmailAnalysisResult) error

	// SaveEmailMultiple は複数案件対応のメール分析結果をデータベースに保存します
	SaveEmailMultiple(ctx context.Context, result *openaidomain.EmailAnalysisMultipleResult) error

	// GetEmailByGmailId はIDでメールを取得します
	GetEmailByGmailId(ctx context.Context, gmail_id string) (*domain.Email, error)

	// EmailExists はメールが既に存在するかチェックします
	EmailExists(ctx context.Context, id string) (bool, error)

	// KeywordExists はキーワードが既に存在するかチェックします
	KeywordExists(word string) (bool, error)

	// PositionExists はポジションが既に存在するかチェックします
	PositionExists(ctx context.Context, word string) (bool, error)

	// WorkTypeExists は業務種別が既に存在するかチェックします
	WorkTypeExists(ctx context.Context, word string) (bool, error)
}
