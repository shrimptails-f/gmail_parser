package di

import (
	"business/tools/gmail"
	"business/tools/gmailService"
	"business/tools/mysql"
	"business/tools/openai"
	"business/tools/oswrapper"

	"go.uber.org/dig"
)

// ProvideCommonDependencies 共通の依存性（例：データベース接続など）を設定する関数
func ProvideCommonDependencies(container *dig.Container, conn *mysql.MySQL, oa *openai.Client, gs *gmailService.Client, gc *gmail.Client, osw *oswrapper.OsWrapper) {
	_ = container.Provide(func() *mysql.MySQL {
		return conn
	})

	_ = container.Provide(func() *openai.Client {
		return oa
	})

	// gmail.Clientの生成に必要
	_ = container.Provide(func() *gmailService.Client {
		return gs
	})

	// GメールのAPI接続クライアント
	_ = container.Provide(func() *gmail.Client {
		return gc
	})

	_ = container.Provide(func() *oswrapper.OsWrapper {
		return osw
	})

	// var wt ct.CustomTime
	// _ = container.Provide(func() ct.WrapperTime {
	// 	return wt
	// })
}

// BuildContainer すべての依存性を統合して設定するコンテナビルダー関数
func BuildContainer(conn *mysql.MySQL, oa *openai.Client, gs *gmailService.Client, gc *gmail.Client, osw *oswrapper.OsWrapper) *dig.Container {
	container := dig.New()

	// 共通の依存性を登録
	ProvideCommonDependencies(container, conn, oa, gs, gc, osw)

	// 各機能群の依存性を登録
	ProvideOpenAiDependencies(container)
	ProvideGmailDependencies(container)
	ProvideEmailStoreDependencies(container)
	ProvidePresentationDependencies(container)

	return container
}
