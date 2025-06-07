package di

import (
	ea "business/internal/emailstore/application"
	ei "business/internal/emailstore/infrastructure"
	"business/tools/mysql"

	"go.uber.org/dig"
)

// ProvideEmailStoreDependencies 解析結果保存を実行する機能群の依存注入設定
func ProvideEmailStoreDependencies(container *dig.Container) {
	// infra
	_ = container.Provide(func(conn *mysql.MySQL) *ei.EmailStoreRepositoryImpl {
		return ei.NewEmailStoreRepository(conn.DB)
	})
	// app
	_ = container.Provide(func(ei *ei.EmailStoreRepositoryImpl) *ea.EmailStoreUseCaseImpl {
		return ea.NewEmailStoreUseCase(ei)
	})
}
