package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	log.Printf("[AUTH] Initializing AuthMiddleware with JWT Secret length: %d", len(jwtSecret))
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// SupabaseAuth 验证 Supabase JWT Token
func (m *AuthMiddleware) SupabaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[AUTH] Processing request: %s %s", c.Request.Method, c.Request.URL.Path)

		// 1. 从请求头获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("[AUTH] Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: 缺少 Authorization 头",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 2. 解析 Bearer Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("[AUTH] Invalid Authorization format (missing 'Bearer ' prefix)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Authorization 格式错误",
				"data":    nil,
			})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Token received (length: %d)", len(tokenString))
		log.Printf("[AUTH] Token first 20 chars: %s...", tokenString[:min(20, len(tokenString))])

		// 3. 验证 JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			log.Printf("[AUTH] JWT Algorithm: %v", token.Method.Alg())
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[AUTH] Invalid signing method: %v", token.Method.Alg())
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			log.Printf("[AUTH] Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Token 无效或已过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("[AUTH] Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Token 无效或已过期",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 4. 提取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[AUTH] Failed to parse token claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: Token Claims 解析失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Token claims: %+v", claims)

		userID, ok := claims["sub"].(string)
		if !ok {
			log.Printf("[AUTH] User ID not found in claims or invalid type. Claims: %+v", claims)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    20001,
				"message": "未授权: 用户ID不存在",
				"data":    nil,
			})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Authentication successful for user: %s", userID)

		// 5. 将用户ID存入上下文
		c.Set("user_id", userID)
		if email, ok := claims["email"].(string); ok {
			c.Set("user_email", email)
		}

		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserIDFromContext 从标准 Context 中获取用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}
