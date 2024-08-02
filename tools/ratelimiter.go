package tools

import (
	"sync"
	"time"
)

// 内存方式解决请求速率方案 用户和权限 多对多关系
var (
	userLimiters = make(map[string]*UserRateLimiter) // 用户速率限制器
	mu           sync.Mutex
)

type UserRateLimiter struct {
	Rules map[string]*RateLimiter
	Mutex sync.Mutex
}

func NewRateLimiter(maxReq int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Requests:    0,
		MaxRequests: maxReq,
		Window:      window,
	}
}

type RateLimiter struct {
	Requests    int           //请求总数
	Timestamp   time.Time     //开始计时时间
	MaxRequests int           //计时时间内最大请求次数
	Window      time.Duration //计时单位
	Mutex       sync.Mutex
}

func GetUserRateLimiter(userId string) (*UserRateLimiter, bool) {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := userLimiters[userId]; exists {
		return limiter, true
	}
	return &UserRateLimiter{
		Rules: make(map[string]*RateLimiter),
	}, false
}

// 增加使用速率的规则
func (u *UserRateLimiter) SetUserRules(item map[string]*RateLimiter) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	for k, v := range item {
		u.Rules[k] = v
	}
}

// 根据用户和路由获取速率限制器
func (ul *UserRateLimiter) GetRateLimiter(path string) *RateLimiter {
	ul.Mutex.Lock()
	defer ul.Mutex.Unlock()
	limiter, ok := ul.Rules[path]
	if !ok {
		return nil
	}
	return limiter
}

func (rl *RateLimiter) AllowRequest() bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	now := time.Now()
	if now.Sub(rl.Timestamp) > rl.Window {
		rl.Timestamp = now
		rl.Requests = 0
	}

	if rl.Requests < rl.MaxRequests {
		rl.Requests++
		return true
	}
	return false
}
