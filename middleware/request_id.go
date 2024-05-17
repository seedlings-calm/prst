package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	//requestId key
	XRequestId = "X-Request-Id"
)

// RequestId 自动增加requestId
func RequestId(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		requestId := c.GetHeader(key)
		if requestId == "" {
			requestId = c.GetHeader(strings.ToLower(key))
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Request.Header.Set(key, requestId)
		c.Set(key, requestId)

		c.Next()
	}
}
