package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/config"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	burst    int           // max burst size
	cleanup  time.Duration // cleanup interval
}

// Visitor represents a rate limit visitor
type Visitor struct {
	tokens     int
	lastSeen   time.Time
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &Visitor{
			tokens:     rl.burst - 1,
			lastSeen:   now,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on time elapsed
	elapsed := now.Sub(visitor.lastRefill)
	tokensToAdd := int(elapsed.Minutes() * float64(rl.rate))

	if tokensToAdd > 0 {
		visitor.tokens += tokensToAdd
		if visitor.tokens > rl.burst {
			visitor.tokens = rl.burst
		}
		visitor.lastRefill = now
	}

	visitor.lastSeen = now

	if visitor.tokens > 0 {
		visitor.tokens--
		return true
	}

	return false
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, visitor := range rl.visitors {
			if now.Sub(visitor.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	rateLimitConfig := config.App.RateLimit
	if !rateLimitConfig.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewRateLimiter(rateLimitConfig.RPS, rateLimitConfig.Burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
