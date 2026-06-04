package httpmw

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipRateLimiter 는 클라이언트 IP 별 토큰 버킷을 관리한다(브루트포스/남용 방지용).
type ipRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*ipBucket
	limit   rate.Limit
	burst   int
	ttl     time.Duration
}

type ipBucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newIPRateLimiter(perMinute, burst int) *ipRateLimiter {
	return &ipRateLimiter{
		buckets: make(map[string]*ipBucket),
		limit:   rate.Every(time.Minute / time.Duration(perMinute)),
		burst:   burst,
		ttl:     10 * time.Minute,
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[ip]
	if !ok {
		b = &ipBucket{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.buckets[ip] = b
	}
	b.lastSeen = time.Now()
	return b.limiter.Allow()
}

// cleanupLoop 는 오래 쓰이지 않은 IP 버킷을 주기적으로 제거해 메모리 증가를 막는다.
func (l *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(l.ttl)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-l.ttl)
		l.mu.Lock()
		for ip, b := range l.buckets {
			if b.lastSeen.Before(cutoff) {
				delete(l.buckets, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimit 은 클라이언트 IP 당 분당 perMinute 회(순간 최대 burst)로 제한하는 미들웨어를 만든다.
// perMinute <= 0 이면 제한을 비활성화한다. 초과 시 429 + {message} 응답.
// 주의: c.ClientIP() 는 신뢰 프록시 설정에 의존하므로, 프록시 뒤에서는 SetTrustedProxies 구성 필요.
func RateLimit(perMinute, burst int) gin.HandlerFunc {
	if perMinute <= 0 {
		return func(c *gin.Context) { c.Next() }
	}
	if burst < 1 {
		burst = 1
	}
	l := newIPRateLimiter(perMinute, burst)
	go l.cleanupLoop()
	return func(c *gin.Context) {
		if !l.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "요청이 너무 많습니다. 잠시 후 다시 시도하세요."})
			return
		}
		c.Next()
	}
}
