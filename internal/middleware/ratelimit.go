package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimit returns a fixed-window limiter keyed by client IP: at most `limit`
// requests per `window`. Single-instance, in-memory — sufficient for a
// single-host deployment; excess requests get 429.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		v, ok := visitors[ip]
		if !ok || now.Sub(v.windowStart) > window {
			// New window for this IP. Opportunistically drop stale entries.
			for k, vis := range visitors {
				if now.Sub(vis.windowStart) > window {
					delete(visitors, k)
				}
			}
			visitors[ip] = &visitor{count: 1, windowStart: now}
			mu.Unlock()
			c.Next()
			return
		}

		v.count++
		if v.count > limit {
			retry := int(window.Seconds() - now.Sub(v.windowStart).Seconds())
			if retry < 0 {
				retry = 0
			}
			mu.Unlock()
			c.Header("Retry-After", strconv.Itoa(retry))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts, please slow down"})
			return
		}
		mu.Unlock()
		c.Next()
	}
}
