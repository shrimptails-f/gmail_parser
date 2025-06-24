package di

import (
	ea "business/internal/emailstore/application"
	ei "business/internal/emailstore/infrastructure"
	ga "business/internal/gmail/application"
	gi "business/internal/gmail/infrastructure"
	gc "business/tools/gmail"
	"business/tools/mysql"

	"go.uber.org/dig"
)

// ProvideGmailDependencies Gmail APIを実行する機能群の依存注入設定
func ProvideGmailDependencies(container *dig.Container) {
	// infra
	_ = container.Provide(func(gc *gc.Client) *gi.GmailConnect {
		return gi.New(gc)
	})
	_ = container.Provide(func(conn *mysql.MySQL) *ei.EmailStoreRepositoryImpl {
		return ei.NewEmailStoreRepository(conn.DB)
	})
	// app
	_ = container.Provide(func(ei *ei.EmailStoreRepositoryImpl) *ea.EmailStoreUseCaseImpl {
		return ea.NewEmailStoreUseCase(ei)
	})
	_ = container.Provide(func(gi *gi.GmailConnect, ea *ea.EmailStoreUseCaseImpl) *ga.GmailUseCase {
		return ga.New(gi, ea)
	})
}
