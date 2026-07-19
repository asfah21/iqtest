package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter menyediakan rate limiting per-IP
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter membuat rate limiter baru
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	// Background cleanup goroutine
	go func() {
		for {
			time.Sleep(window)
			rl.mu.Lock()
			now := time.Now()
			for ip, times := range rl.requests {
				// Filter out old entries
				var recent []time.Time
				for _, t := range times {
					if now.Sub(t) <= window {
						recent = append(recent, t)
					}
				}
				if len(recent) == 0 {
					delete(rl.requests, ip)
				} else {
					rl.requests[ip] = recent
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Clean old entries for this IP
	var recent []time.Time
	for _, t := range rl.requests[ip] {
		if now.Sub(t) <= rl.window {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rl.limit {
		rl.requests[ip] = recent
		return false
	}

	rl.requests[ip] = append(recent, now)
	return true
}

// RateLimitMiddleware returns a Gin middleware function for rate limiting
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan. Silakan coba lagi nanti.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
