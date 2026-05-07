# 🚨 URLSentinel

A blazing-fast, concurrent URL health checker written in Go. Monitor your web services, APIs, and external links with ease.

[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org/doc/devel/release)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Use Cases](#use-cases)
- [Configuration](#configuration)
- [Output](#output)
- [Examples](#examples)
- [Contributing](#contributing)

## ✨ Features

- **⚡ Concurrent Checking**: Adjust worker count for optimal performance
- **⏱️ Configurable Timeouts**: Set per-request timeout durations
- **📊 Detailed Reports**: JSON export for integration with monitoring tools
- **🔄 Context-Aware**: Graceful cancellation and cleanup
- **📄 Batch Processing**: Check hundreds of URLs from a file
- **🎨 Rich CLI Output**: Visual status indicators and real-time feedback
- **🏎️ High Performance**: Processes 1000+ URLs efficiently using goroutines

## 🚀 Installation

### Prerequisites
- Go 1.20 or higher

### From Source

```bash
git clone https://github.com/ammarkhan575/url-sentinel.git
cd url-sentinel
go build
```

### Quick Install with Go

```bash
go install github.com/ammarkhan575/url-sentinel@latest
```

## 📖 Usage

### Basic Command

```bash
./url-sentinel -file urls.txt -workers 10 -timeout 5s
```

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-file` | **required** | Path to file containing URLs (one per line) |
| `-workers` | 10 | Number of concurrent workers |
| `-timeout` | 5s | HTTP request timeout duration (e.g., 3s, 100ms) |
| `-output` | empty | Path to save JSON report (optional) |

## 🎯 Use Cases

### 1. **Website Health Monitoring**
Monitor your company's public-facing websites and ensure they're always available.

```bash
./url-sentinel -file websites.txt -workers 5 -timeout 3s -output report.json
```

### 2. **API Endpoint Verification**
Verify that your microservices API endpoints are responding correctly.

```bash
./url-sentinel -file api-endpoints.txt -workers 20 -timeout 2s
```

### 3. **Third-Party Service Availability Check**
Check if third-party services your application depends on are online.

```bash
./url-sentinel -file external-services.txt -workers 15 -timeout 10s
```

### 4. **Link Validation in CI/CD Pipeline**
Validate all links in your documentation or website before deployment.

```bash
./url-sentinel -file generated-links.txt -workers 50 -timeout 5s -output ci-report.json
```

### 5. **Load Balancer Health Verification**
Check multiple instances behind a load balancer to ensure distribution.

```bash
./url-sentinel -file lb-instances.txt -workers 8 -timeout 2s
```

### 6. **Scheduled Health Checks (Cron)**
Run periodic checks and log results for historical analysis.

```bash
0 */6 * * * /path/to/url-sentinel -file urls.txt -output /logs/health-$(date +\%Y\%m\%d-\%H\%M\%S).json
```

### 7. **Dependency Service Monitoring**
Before deploying, verify all external dependencies are accessible.

```bash
./url-sentinel -file dependencies.txt -workers 10 -timeout 5s
```

### 8. **Content Delivery Network (CDN) Edge Validation**
Check CDN edge locations to ensure content is being served globally.

```bash
./url-sentinel -file cdn-edges.txt -workers 30 -timeout 3s -output cdn-status.json
```

### 9. **Webhook Integration Testing**
Verify webhook endpoints are accepting connections before triggering events.

```bash
./url-sentinel -file webhooks.txt -workers 5 -timeout 2s
```

### 10. **Disaster Recovery Drills**
Test backup and failover systems are accessible and responding.

```bash
./url-sentinel -file backup-systems.txt -workers 10 -timeout 10s -output dr-drill-results.json
```

## ⚙️ Configuration

### Input File Format

Create a text file with one URL per line:

```
https://example.com
https://api.example.com/health
https://cdn.example.com
http://internal-service:8080
https://monitoring.example.com/status
```

### Optimal Worker Count

- **Development/Testing**: 5-10 workers
- **Small Batch (< 100 URLs)**: 10-15 workers
- **Medium Batch (100-500 URLs)**: 20-30 workers
- **Large Batch (500+ URLs)**: 50-100+ workers

### Timeout Recommendations

- **Local Services**: 1-2 seconds
- **Public APIs**: 3-5 seconds
- **Slow/Heavy Endpoints**: 10-30 seconds
- **CDN Checks**: 5-10 seconds

## 📊 Output

### Console Output

```
Welcome to URLSentinel v0.1
Checking 6 URLs with 5 workers, timeout 3s...
✓ https://github.com                                 up 272ms
✗ https://this-does-not-exist-xyz.com                down 395ms
✓ https://cloudflare.com                             up 729ms
✓ https://google.com                                 up 841ms
✗ https://httpstat.us/500                            down 3001ms
✗ https://example.org                                timeout 3001ms

Done: 3 up, 3 down, 1 timeout. Time: 3.001s
```

### JSON Report

Save detailed results with `-output report.json`:

```json
{
  "total": 6,
  "up": 3,
  "down": 2,
  "timeouts": 1,
  "total_time_ms": 3001000000,
  "results": [
    {
      "url": "https://github.com",
      "status": "up",
      "status_code": 200,
      "latency_ms": 272,
      "checked_at": "2026-05-07T15:46:09.844609+05:30"
    },
    {
      "url": "https://this-does-not-exist-xyz.com",
      "status": "down",
      "latency_ms": 395,
      "error": "dial tcp: lookup this-does-not-exist-xyz.com: no such host",
      "checked_at": "2026-05-07T15:46:09.572577+05:30"
    }
  ]
}
```

## 💡 Examples

### Example 1: Monitor Your Production Stack

```bash
# Create urls.txt
cat > urls.txt << EOF
https://app.example.com
https://api.example.com/v1/health
https://cdn.example.com
https://db-backup.example.com
EOF

# Run the checker
./url-sentinel -file urls.txt -workers 10 -timeout 5s -output prod-health.json
```

### Example 2: CI/CD Integration

```bash
#!/bin/bash
# Validate all dependencies before deployment

./url-sentinel -file dependencies.txt -workers 20 -timeout 3s -output ci-deps.json

if [ $? -eq 0 ]; then
  echo "✓ All dependencies are healthy"
  cat ci-deps.json | jq '.up, .down, .timeouts'
else
  echo "✗ Some services are down"
  exit 1
fi
```

### Example 3: Automated Monitoring Script

```bash
#!/bin/bash
# Check every 6 hours

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
LOG_DIR="/var/log/url-sentinel"
mkdir -p "$LOG_DIR"

./url-sentinel \
  -file /etc/url-sentinel/urls.txt \
  -workers 30 \
  -timeout 10s \
  -output "$LOG_DIR/check-$TIMESTAMP.json"

# Send alert if any URL is down
DOWN=$(jq '.down' "$LOG_DIR/check-$TIMESTAMP.json")
if [ "$DOWN" -gt 0 ]; then
  echo "Alert: $DOWN URLs are down" | mail -s "URLSentinel Alert" ops@example.com
fi
```

### Example 4: Performance Testing

```bash
# Test with 1000 URLs and measure performance
time ./url-sentinel -file large-urls.txt -workers 100 -timeout 15s -output results.json
```

## 🏗️ Architecture

URLSentinel uses a **worker pool pattern** with goroutines:

1. **Job Queue**: URLs are buffered in a channel
2. **Worker Goroutines**: Process URLs concurrently
3. **Result Collection**: Results are gathered as workers complete
4. **Context Cancellation**: Clean shutdown on timeout or error

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with race detection
go run -race . -file urls.txt -workers 5 -timeout 3s
```

## 📝 Status Codes

| Status | Meaning |
|--------|---------|
| `up` | URL responded with 2xx or 3xx status code |
| `down` | URL returned 4xx/5xx, DNS failed, or connection refused |
| `timeout` | Request exceeded the specified timeout duration |

## 🐛 Troubleshooting

### "Error reading file: no such file or directory"
Ensure the `-file` flag points to an existing file with valid URLs.

### Timeouts occurring frequently
Increase the `-timeout` value or reduce `-workers` to avoid overwhelming the target.

### High memory usage with large URL lists
Reduce `-workers` or process URLs in batches.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙋 Support

Found a bug? Have a feature request? [Open an issue](https://github.com/ammarkhan575/url-sentinel/issues)

---

**Built with ❤️ using Go**
