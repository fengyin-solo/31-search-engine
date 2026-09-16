package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"searchengine/pkg/httpx"
)

// authMiddleware 校验 X-Auth-Token 请求头（仅对 /api/ 前缀接口生效，静态页面不鉴权）。
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		if s.cfg == nil || s.cfg.AuthToken == "" {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-Auth-Token") != s.cfg.AuthToken {
			httpx.Unauthorized(w, "无效或缺失的访问令牌")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// tokenBucket 令牌桶。
type tokenBucket struct {
	tokens float64
	last   time.Time
}

// rateLimiter 按客户端 IP 限流的令牌桶限流器。
type rateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64
	capacity float64
}

func newRateLimiter(ratePerMin int) *rateLimiter {
	rate := float64(ratePerMin) / 60.0
	if rate <= 0 {
		rate = 10
	}
	return &rateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     rate,
		capacity: rate * 2,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	b, ok := rl.buckets[key]
	if !ok {
		b = &tokenBucket{tokens: rl.capacity, last: now}
		rl.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * rl.rate
	}
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// rateLimitMiddleware 按 IP 限流。
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	rl := newRateLimiter(s.cfg.RateLimit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			return
		}
		next.ServeHTTP(w, r)
	})
}
