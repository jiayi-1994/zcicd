package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit applies a global token-bucket rate limiter.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    42900,
				"message": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimit applies a per-IP token-bucket rate limiter.
// Stale entries are cleaned up every 3 minutes.
func IPRateLimit(rps float64, burst int) gin.HandlerFunc {
	var clients sync.Map

	// Background cleanup of stale entries.
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			clients.Range(func(key, value any) bool {
				e := value.(*ipEntry)
				if time.Since(e.lastSeen) > 5*time.Minute {
					clients.Delete(key)
				}
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, ok := clients.Load(ip)
		var e *ipEntry
		if ok {
			e = val.(*ipEntry)
			e.lastSeen = now
		} else {
			e = &ipEntry{
				limiter:  rate.NewLimiter(rate.Limit(rps), burst),
				lastSeen: now,
			}
			actual, loaded := clients.LoadOrStore(ip, e)
			if loaded {
				e = actual.(*ipEntry)
				e.lastSeen = now
			}
		}

		if !e.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    42900,
				"message": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
