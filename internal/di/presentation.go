package di

import (
	"business/internal/app/presentation"
	ea "business/internal/emailstore/application"
	ga "business/internal/gmail/application"
	aiapp "business/internal/openAi/application"

	"go.uber.org/dig"
)

// ProvidePresentationDependencies プレゼンテーション層の依存注入設定
func ProvidePresentationDependencies(container *dig.Container) {
	// AnalyzeEmailControllerの依存注入
	_ = container.Provide(func(
		ea *ea.EmailStoreUseCaseImpl,
		ga *ga.GmailUseCase,
		aiapp *aiapp.UseCase,
	) *presentation.AnalyzeEmailController {
		return presentation.NewAnalyzeEmailController(ea, ga, aiapp)
	})
}
