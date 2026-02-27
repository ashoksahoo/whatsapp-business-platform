package middleware

import (
	"sync"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerMinute) / 60, // Convert to per second
		burst:    requestsPerMinute / 10,             // 10% burst capacity
	}
}

// getLimiter gets or creates a limiter for a key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// RateLimitMiddleware applies rate limiting per API key
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute)

	return func(c *gin.Context) {
		// Get API key ID from context (set by auth middleware)
		apiKeyID, exists := c.Get("api_key_id")
		if !exists {
			// If not authenticated, use IP as key
			apiKeyID = c.ClientIP()
		}

		key := apiKeyID.(string)
		rateLimiter := limiter.getLimiter(key)

		if !rateLimiter.Allow() {
			// Set rate limit headers
			c.Header("X-RateLimit-Limit", string(rune(requestsPerMinute)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", string(rune(time.Now().Add(time.Minute).Unix())))

			utils.ErrorJSON(c, errors.NewRateLimitError())
			c.Abort()
			return
		}

		c.Next()
	}
}
