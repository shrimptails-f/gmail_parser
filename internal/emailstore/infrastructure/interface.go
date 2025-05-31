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

	// GetEmailByID はIDでメールを取得します
	GetEmailByID(ctx context.Context, id string) (*domain.Email, error)

	// EmailExists はメールが既に存在するかチェックします
	EmailExists(ctx context.Context, id string) (bool, error)
}
