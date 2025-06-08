package di

import (
	ga "business/internal/gmail/application"
	gi "business/internal/gmail/infrastructure"
	gc "business/tools/gmail"

	"go.uber.org/dig"
)

// ProvideGmailDependencies Gmail APIを実行する機能群の依存注入設定
func ProvideGmailDependencies(container *dig.Container) {
	// infra
	_ = container.Provide(func(gc *gc.Client) *gi.GmailConnect {
		return gi.New(gc)
	})
	// app
	_ = container.Provide(func(gi *gi.GmailConnect) *ga.GmailUseCase {
		return ga.New(gi)
	})
}
