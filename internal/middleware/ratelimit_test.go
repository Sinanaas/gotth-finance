package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newLimitedRouter(limit int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimit(limit, window))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func hit(r *gin.Engine, ip string) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":40000"
	r.ServeHTTP(w, req)
	return w.Code
}

func TestRateLimit_BlocksOverLimit(t *testing.T) {
	r := newLimitedRouter(3, time.Minute)
	for i := 1; i <= 3; i++ {
		if code := hit(r, "1.2.3.4"); code != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i, code)
		}
	}
	if code := hit(r, "1.2.3.4"); code != http.StatusTooManyRequests {
		t.Fatalf("4th request: got %d, want 429", code)
	}
}

func TestRateLimit_WindowResets(t *testing.T) {
	r := newLimitedRouter(1, 20*time.Millisecond)
	if hit(r, "5.6.7.8") != http.StatusOK {
		t.Fatal("first request should pass")
	}
	if hit(r, "5.6.7.8") != http.StatusTooManyRequests {
		t.Fatal("second request within window should be blocked")
	}
	time.Sleep(30 * time.Millisecond)
	if hit(r, "5.6.7.8") != http.StatusOK {
		t.Fatal("request after window should pass")
	}
}

func TestRateLimit_PerIP(t *testing.T) {
	r := newLimitedRouter(1, time.Minute)
	if hit(r, "9.9.9.9") != http.StatusOK {
		t.Fatal("IP A first request should pass")
	}
	if hit(r, "8.8.8.8") != http.StatusOK {
		t.Fatal("IP B should be tracked independently")
	}
}
