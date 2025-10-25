package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			slog.Error("Request error",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": err.Error(),
			})
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
