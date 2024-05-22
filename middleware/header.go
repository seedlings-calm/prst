package middleware

import (
	"github.com/gin-gonic/gin"
)

func Options() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Origin, Content-Type, Accept,X-Request-Id,X-Forwarded-For,X-Real-IP")
			c.Header("Access-Control-Allow-Credentials", "true")
			// c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(200)
		}
	}
}

// 设置http header安全
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	// c.Header("Referrer-Policy", "same-origin")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
	// Also consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "default-src 'self'")
}
