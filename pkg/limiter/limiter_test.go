package limiter

import (
	"testing"
	"time"
)

func TestTokenBucketBurst(t *testing.T) {
	tb := NewTokenBucket(10, 5) // 10 tokens/sec, burst 5

	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Fatalf("expected token %d to be allowed", i)
		}
	}

	// 6th immediate request should be rejected
	if tb.Allow() {
		t.Fatalf("expected 6th burst request to be rejected")
	}

	// Wait 250ms -> should refill ~2.5 tokens
	time.Sleep(250 * time.Millisecond)
	if !tb.Allow() {
		t.Fatalf("expected request after refill to be allowed")
	}
}

func TestIPLimiterIsolation(t *testing.T) {
	il := NewIPLimiter(5, 2)

	if !il.AllowIP("192.168.1.1") || !il.AllowIP("192.168.1.1") {
		t.Fatalf("IP 1 burst failed")
	}
	if il.AllowIP("192.168.1.1") {
		t.Fatalf("IP 1 should be rate-limited")
	}

	// IP 2 must not be affected by IP 1 limit
	if !il.AllowIP("10.0.0.1") {
		t.Fatalf("IP 2 should be allowed independent of IP 1")
	}
}
