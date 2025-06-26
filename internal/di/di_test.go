package di

import (
	"business/internal/app/presentation"
	"business/tools/gmail"
	"business/tools/gmailService"
	"business/tools/mysql"
	"business/tools/openai"
	"business/tools/oswrapper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildContainer_NoError(t *testing.T) {
	// ダミー（空実装）具象を生成
	conn := &mysql.MySQL{}
	oa := &openai.Client{}
	gs := &gmailService.Client{}
	gc := &gmail.Client{}
	osw := &oswrapper.OsWrapper{}

	container := BuildContainer(conn, oa, gs, gc, osw)

	// invokeだけを行い、実行はしない（副作用なし）
	err := container.Invoke(func(
		_ *mysql.MySQL,
		_ *openai.Client,
		_ *gmailService.Client,
		_ *oswrapper.OsWrapper,
	) {
		// 何もしない
	})

	assert.NoError(t, err)
}

func TestBuildContainer_WithPresentationLayer(t *testing.T) {
	// ダミー（空実装）具象を生成
	conn := &mysql.MySQL{}
	oa := &openai.Client{}
	gs := &gmailService.Client{}
	gc := &gmail.Client{}
	osw := &oswrapper.OsWrapper{}

	container := BuildContainer(conn, oa, gs, gc, osw)

	// presentation層の依存注入をテスト
	err := container.Invoke(func(controller *presentation.AnalyzeEmailController) {
		// controllerが正常に注入されることを確認
		assert.NotNil(t, controller)
	})

	assert.NoError(t, err)
}
