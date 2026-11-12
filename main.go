package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IsmeetPandey/goguard-proxy/pkg/proxy"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	target := flag.String("target", "http://localhost:3000", "Upstream target server URL")
	rate := flag.Float64("rate", 50.0, "Allowed requests per second per IP")
	burst := flag.Float64("burst", 25.0, "Maximum burst capacity")
	flag.Parse()

	fmt.Println("============================================================")
	fmt.Println("🛡️  GOGUARD-PROXY v1.0.0")
	fmt.Println("├─ High-Performance Rate-Limiting Reverse Proxy")
	fmt.Println("└─ Author: Ismeet Pandey")
	fmt.Println("============================================================")

	gp, err := proxy.NewGuardProxy(proxy.Config{
		TargetURL:  *target,
		RatePerSec: *rate,
		BurstLimit: *burst,
	})
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      gp,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("🚀 GoGuard listening on http://0.0.0.0:%d -> forwarding to %s", *port, *target)
		log.Printf("📊 Metrics available at http://0.0.0.0:%d/_goguard/metrics", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down GoGuard...")
}
