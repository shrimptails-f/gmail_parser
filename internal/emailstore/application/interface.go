// Package application はメール保存機能のアプリケーション層を提供します。
// このファイルはメール保存に関するインターフェースを定義します。
package application

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

// EmailStoreUseCase はメール保存のユースケースインターフェースです
type EmailStoreUseCase interface {
	// SaveEmailAnalysisResult はメール分析結果を保存します
	SaveEmailAnalysisResult(ctx context.Context, result *openaidomain.EmailAnalysisResult) error
}
