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

		// Log incoming request with Authorization header (for debugging)
		authHeader := c.GetHeader("Authorization")
		hasAuth := authHeader != ""
		slog.Info("Incoming request",
			"method", method,
			"path", path,
			"has_auth", hasAuth,
		)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logFields := []any{
			"method", method,
			"path", path,
			"status", statusCode,
			"duration", duration.String(),
			"client_ip", clientIP,
		}

		// 添加错误信息（如果有）
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			slog.Error("HTTP Request", logFields...)
		} else if statusCode >= 400 {
			slog.Warn("HTTP Request", logFields...)
		} else {
			slog.Info("HTTP Request", logFields...)
		}
	}
}
