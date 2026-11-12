package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/IsmeetPandey/goguard-proxy/pkg/limiter"
)

type Config struct {
	TargetURL    string
	RatePerSec   float64
	BurstLimit   float64
	HeaderPrefix string
}

type GuardProxy struct {
	proxy       *httputil.ReverseProxy
	limiter     *limiter.IPLimiter
	totalReqs   atomic.Uint64
	blockedReqs atomic.Uint64
}

func NewGuardProxy(cfg Config) (*GuardProxy, error) {
	target, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-By", "GoGuard-Proxy/v1.0")
		req.Header.Set("X-Proxy-Arrival-Time", time.Now().UTC().Format(time.RFC3339Nano))
	}

	return &GuardProxy{
		proxy:   reverseProxy,
		limiter: limiter.NewIPLimiter(cfg.RatePerSec, cfg.BurstLimit),
	}, nil
}

func (gp *GuardProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gp.totalReqs.Add(1)

	// Metrics endpoint
	if r.URL.Path == "/_goguard/metrics" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(gp.Metrics()))
		return
	}

	clientIP := extractIP(r)
	if !gp.limiter.AllowIP(clientIP) {
		gp.blockedReqs.Add(1)
		w.Header().Set("Retry-After", "1")
		http.Error(w, "429 Too Many Requests - GoGuard Rate Limit Exceeded", http.StatusTooManyRequests)
		log.Printf("[RATE_LIMIT] Blocked client %s on %s %s", clientIP, r.Method, r.URL.Path)
		return
	}

	gp.proxy.ServeHTTP(w, r)
}

func (gp *GuardProxy) Metrics() string {
	return strings.Join([]string{
		"# HELP goguard_requests_total Total number of incoming proxy requests",
		"# TYPE goguard_requests_total counter",
		"goguard_requests_total " + stringFromUint(gp.totalReqs.Load()),
		"# HELP goguard_requests_blocked_total Total number of rate-limited requests",
		"# TYPE goguard_requests_blocked_total counter",
		"goguard_requests_blocked_total " + stringFromUint(gp.blockedReqs.Load()),
		"",
	}, "
")
}

func extractIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func stringFromUint(n uint64) string {
	return strings.TrimRight(strings.TrimSpace(string(rune('0'+n%10))), "")
}
