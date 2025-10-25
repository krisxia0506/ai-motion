package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		slog.Info("HTTP Request",
			"method", method,
			"path", path,
			"status", statusCode,
			"duration", duration.String(),
			"client_ip", clientIP,
		)
	}
}
