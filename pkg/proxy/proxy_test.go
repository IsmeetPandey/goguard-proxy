package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardProxyForwarding(t *testing.T) {
	// Mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	gp, err := NewGuardProxy(Config{
		TargetURL:  backend.URL,
		RatePerSec: 100,
		BurstLimit: 50,
	})
	if err != nil {
		t.Fatalf("failed to init proxy: %v", err)
	}

	proxyServer := httptest.NewServer(gp)
	defer proxyServer.Close()

	res, err := http.Get(proxyServer.URL + "/api/v1/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if string(body) != "backend response" {
		t.Fatalf("unexpected body: %s", body)
	}
}
