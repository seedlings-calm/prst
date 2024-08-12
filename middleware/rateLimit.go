package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/pkg/ratelimiter"
)

var (
	mu     sync.Mutex
	Router = map[string]*RateType{
		"/api/v1/user/collect":                   {MaxRequests: 20, Is: true, Window: 1 * time.Minute},
		"/api/v1/index/getCollect":               {MaxRequests: 20, Is: true, Window: 1 * time.Minute},
		"/api/v1/user/practtioner/getQuestion/1": {MaxRequests: 3, Is: true, Window: 24 * 60 * time.Minute},
		"/api/v1/user/practtioner/getQuestion/2": {MaxRequests: 3, Is: true, Window: 24 * 60 * time.Minute},
		"/api/v1/user/practtioner/getQuestion/3": {MaxRequests: 3, Is: true, Window: 24 * 60 * time.Minute},
	}
)

type RateType struct {
	MaxRequests int           //限制次数
	Is          bool          //是否启用路由限制
	Window      time.Duration //限时时间
}

// 请求速率限制
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		userId := c.GetString("userId")
		//如果路由不存在检测里面，直接跳过
		val, ok := Router[path]
		if !ok || !val.Is || userId == "" {
			c.Next()
			return
		}
		mu.Lock()
		defer mu.Unlock()
		userLimiter, ok := ratelimiter.GetUserRateLimiter(userId)
		if !ok {
			rates := make(map[string]*ratelimiter.RateLimiter)
			for k, v := range Router {
				item := ratelimiter.NewRateLimiter(v.MaxRequests, v.Window)
				rates[k] = item
			}
			userLimiter.SetUserRules(userId, rates)
		}
		limiter := userLimiter.GetRateLimiter(path)
		if !limiter.AllowRequest() {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "too manys request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
