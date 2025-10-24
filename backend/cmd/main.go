package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "ai-motion",
		})
	})

	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// 小说管理
		novel := v1.Group("/novel")
		{
			novel.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "upload novel endpoint"})
			})
			novel.POST("/:id/parse", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "parse novel endpoint"})
			})
		}

		// 角色管理
		v1.GET("/characters/:novel_id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "get characters endpoint"})
		})

		// 生成服务
		generate := v1.Group("/generate")
		{
			generate.POST("/scene", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "generate scene endpoint"})
			})
			generate.POST("/voice", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "generate voice endpoint"})
			})
		}

		// 动漫导出
		v1.POST("/anime/:id/export", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "export anime endpoint"})
		})
	}

	log.Println("Starting AI-Motion server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
