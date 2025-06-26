package v1

import (
	"business/internal/app/presentation"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func NewRouter(g *gin.Engine, container *dig.Container) *gin.Engine {
	ctx := context.Background()
	g.GET("/", func(c *gin.Context) {
		// c.Status(http.StatusNoContent)

		// c.Status(http.StatusBadRequest)
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	g.GET("/openAi-email-analysis", func(c *gin.Context) {
		err := container.Invoke(func(p *presentation.AnalyzeEmailController) {
			err := p.SaveEmailAnalysisResult(ctx)
			if err != nil {
				fmt.Printf("Eメール解析エラー: %v", err)
				return
			}
		})
		if err != nil {
			fmt.Printf("Eメール解析エラー: %v", err)
			return
		}
		c.Status(http.StatusNoContent)
	})

	return g
}
