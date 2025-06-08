package di

import (
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
		_ *gmail.Client,
		_ *oswrapper.OsWrapper,
	) {
		// 何もしない
	})

	assert.NoError(t, err)
}
