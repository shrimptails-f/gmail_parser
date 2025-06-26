package v1

import (
	"business/internal/app/presentation"
	"context"
	"fmt"
	"net/http"
	"strings"

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
		var innerErr error
		err := container.Invoke(func(p *presentation.AnalyzeEmailController) {
			innerErr = p.SaveEmailAnalysisResult(c, ctx)
		})
		if innerErr != nil {
			if strings.Contains(innerErr.Error(), "BadRequest") {
				c.Status(http.StatusBadRequest)
				return
			}
			fmt.Printf("Eメール解析エラー: %v", innerErr)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err != nil {
			fmt.Printf("Eメール解析エラー: %v", err)
			c.Status(http.StatusInternalServerError)
		}
		c.Status(http.StatusOK)
	})

	return g
}
