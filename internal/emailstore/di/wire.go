// Package di はメール保存機能の依存性注入を提供します。
// このファイルはメール保存に関する依存性の設定を定義します。
package di

import (
	"business/internal/emailstore/application"
	"business/internal/emailstore/infrastructure"

	"gorm.io/gorm"
)

// ProvideEmailStoreDependencies はメール保存機能の依存性を提供します
func ProvideEmailStoreDependencies(db *gorm.DB) application.EmailStoreUseCase {
	emailStoreRepository := infrastructure.NewEmailStoreRepository(db)
	emailStoreUseCase := application.NewEmailStoreUseCase(emailStoreRepository)
	return emailStoreUseCase
}
