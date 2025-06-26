package di

import (
	ea "business/internal/emailstore/application"
	ei "business/internal/emailstore/infrastructure"
	ga "business/internal/gmail/application"
	gi "business/internal/gmail/infrastructure"
	gc "business/tools/gmail"
	gs "business/tools/gmailService"
	"business/tools/mysql"
	"business/tools/oswrapper"

	"go.uber.org/dig"
)

// ProvideGmailDependencies Gmail APIを実行する機能群の依存注入設定
func ProvideGmailDependencies(container *dig.Container) {
	// infra - GmailConnectはgmailService.ClientInterfaceを使用するように修正が必要
	_ = container.Provide(func(gs *gs.Client, gc *gc.Client, osw *oswrapper.OsWrapper) *gi.GmailConnect {
		return gi.New(gs, gc, osw)
	})
	_ = container.Provide(func(conn *mysql.MySQL) *ei.Repository {
		return ei.New(conn.DB)
	})
	// app
	_ = container.Provide(func(ei *ei.Repository) *ea.UseCase {
		return ea.New(ei)
	})
	_ = container.Provide(func(gi *gi.GmailConnect, ea *ea.UseCase) *ga.GmailUseCase {
		return ga.New(gi, ea)
	})
}
