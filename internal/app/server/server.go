package server

import (
	v1 "business/internal/app/router"
	"business/internal/di"
	"business/tools/gmail"
	"business/tools/gmailService"
	"business/tools/mysql"
	"business/tools/openai"
	"business/tools/oswrapper"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run() {
	g := gin.Default()

	osw := oswrapper.New()

	// DBインスタンス生成
	db, err := mysql.New()
	if err != nil {
		fmt.Printf("DB 初期化時にエラーが発生しました。:%v \n,", err)
		return
	}

	// OpenAiクライアント作成
	apiKey := osw.GetEnv("OPENAI_API_KEY")
	oa := openai.New(apiKey)
	gs := gmailService.New()
	gc := gmail.New()

	// DIを行う
	container := di.BuildContainer(db, oa, gs, gc, osw)

	router := v1.NewRouter(g, container)
	router.Run(":8080")
}
