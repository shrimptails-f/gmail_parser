package di

import (
	aiapp "business/internal/openAi/application"
	aiinfra "business/internal/openAi/infrastructure"
	"business/tools/openai"
	"business/tools/oswrapper"

	"go.uber.org/dig"
)

// ProvideOpenAiDependencies OpenAi APIを実行する機能群の依存注入設定
func ProvideOpenAiDependencies(container *dig.Container) {
	// infra
	_ = container.Provide(func(oa *openai.Client) *aiinfra.Analyzer {
		return aiinfra.NewAnalyzer(oa)
	})
	// app
	_ = container.Provide(func(r *aiinfra.Analyzer, osw *oswrapper.OsWrapper) *aiapp.UseCase {
		return aiapp.NewUseCase(r, osw)
	})
}
