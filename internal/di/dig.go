package di

import (
	"business/tools/mysql"
	"business/tools/openai"
	"business/tools/oswrapper"

	"go.uber.org/dig"
)

// ProvideCommonDependencies 共通の依存性（例：データベース接続など）を設定する関数
func ProvideCommonDependencies(container *dig.Container, conn *mysql.MySQL, oa *openai.Client, osw *oswrapper.OsWrapper) {
	_ = container.Provide(func() *mysql.MySQL {
		return conn
	})

	_ = container.Provide(func() *openai.Client {
		return oa
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
func BuildContainer(conn *mysql.MySQL, oa *openai.Client, osw *oswrapper.OsWrapper) *dig.Container {
	container := dig.New()

	// 共通の依存性を登録
	ProvideCommonDependencies(container, conn, oa, osw)

	// 各機能群の依存性を登録
	ProvideOpenAiDependencies(container)

	return container
}
