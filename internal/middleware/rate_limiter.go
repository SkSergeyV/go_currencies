package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	limiters   map[string]*rate.Limiter
	mu         sync.RWMutex
	r          rate.Limit
	b          int
	trustedIPs map[string]struct{}
}

func NewIPRateLimiter(limit int, duration time.Duration, trustedIPs []string) *IPRateLimiter {
	l := &IPRateLimiter{
		limiters:   make(map[string]*rate.Limiter),
		r:          rate.Every(duration / time.Duration(limit)),
		b:          limit,
		trustedIPs: make(map[string]struct{}, len(trustedIPs)),
	}

	for _, ip := range trustedIPs {
		l.trustedIPs[ip] = struct{}{}
	}
	return l
}

func (l *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.RLock()
	limiter, ok := l.limiters[ip]
	l.mu.RUnlock()

	if ok {
		return limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if limiter, ok = l.limiters[ip]; ok {
		return limiter
	}
	limiter = rate.NewLimiter(l.r, l.b)
	l.limiters[ip] = limiter
	return limiter
}

func RateLimiter(limiter *IPRateLimiter) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()

		//--- WHITELIST: проверка  IP ---
		if _, ok := limiter.trustedIPs[ip]; ok {
			return c.Next()
		}
		// -----------------------------------------

		if !limiter.getLimiter(ip).Allow() {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
		}
		return c.Next()
	}
}
