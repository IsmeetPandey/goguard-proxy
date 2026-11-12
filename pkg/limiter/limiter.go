package limiter

import (
	"sync"
	"time"
)

// TokenBucket implements a thread-safe token bucket rate limiter.
type TokenBucket struct {
	rate       float64 // tokens per second
	capacity   float64 // maximum burst capacity
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter.
func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

// Allow checks if 1 token is available, consumes it, and returns true.
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1.0)
}

// AllowN checks if n tokens are available and consumes them.
func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	// Refill tokens according to elapsed duration
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

// IPLimiter manages per-IP rate limiting instances.
type IPLimiter struct {
	rate     float64
	capacity float64
	buckets  map[string]*TokenBucket
	mu       sync.RWMutex
}

// NewIPLimiter creates a per-IP rate limiter registry.
func NewIPLimiter(rate, capacity float64) *IPLimiter {
	return &IPLimiter{
		rate:     rate,
		capacity: capacity,
		buckets:  make(map[string]*TokenBucket),
	}
}

// AllowIP checks if an IP is allowed to execute a request.
func (il *IPLimiter) AllowIP(ip string) bool {
	il.mu.RLock()
	bucket, exists := il.buckets[ip]
	il.mu.RUnlock()

	if !exists {
		il.mu.Lock()
		// Double-check after acquiring write lock
		bucket, exists = il.buckets[ip]
		if !exists {
			bucket = NewTokenBucket(il.rate, il.capacity)
			il.buckets[ip] = bucket
		}
		il.mu.Unlock()
	}

	return bucket.Allow()
}
