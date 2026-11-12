# 🛡️ goguard-proxy

<div align="center">

[![Go CI](https://github.com/IsmeetPandey/goguard-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/IsmeetPandey/goguard-proxy/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/IsmeetPandey/goguard-proxy)](https://goreportcard.com/report/github.com/IsmeetPandey/goguard-proxy)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Language: Go](https://img.shields.io/badge/Language-Go-00ADD8.svg)](https://go.dev/)

**High-performance HTTP reverse proxy with adaptive token-bucket rate limiting, automatic client isolation, and Prometheus metrics.**

</div>

---

## ⚡ Highlights

- **🔒 Atomic Token-Bucket Limiter**: Lock-minimized per-IP rate limiting supporting fractional token refills.
- **🚀 Ultra-Low Latency Overhead**: Sub-millisecond forwarding delay with zero runtime heap allocations on fast paths.
- **📊 Native Prometheus Metrics**: Real-time traffic inspection at `/_goguard/metrics`.
- **🛡️ Header Injection**: Transparently attaches `X-Forwarded-By` and arrival timestamps for distributed tracing.

---

## 📦 Getting Started

### Run with Go
```bash
git clone https://github.com/IsmeetPandey/goguard-proxy.git
cd goguard-proxy
go build -o goguard main.go
./goguard -port 8080 -target http://localhost:3000 -rate 50 -burst 25
```

---

## 🛠️ CLI Flags

| Flag | Default | Description |
|---|---|---|
| `-port` | `8080` | Port to bind GoGuard on |
| `-target` | `http://localhost:3000` | Upstream destination URL |
| `-rate` | `50.0` | Sustained requests per second allowed per IP |
| `-burst` | `25.0` | Maximum burst requests allowed above rate |

---

## 📄 License

MIT License © 2026 [Ismeet Pandey](https://github.com/IsmeetPandey).
