# vanity-go

A lightweight Go vanity import path server that allows you to use custom domain names for your Go packages.

## What is a Go Vanity Import Path?

Go vanity import paths allow you to use your own domain for Go package imports instead of directly referencing the hosting service. For example:

```go
import "go.gllm.dev/mypackage"
```

Instead of:

```go
import "github.com/gllm-dev/mypackage"
```

## Features

- 🚀 **Simple and lightweight** - Minimal dependencies and fast startup
- 🐳 **Docker-ready** - Multi-stage Dockerfile for tiny production images
- 🔧 **Easy configuration** - Just two environment variables
- 📦 **Works with any Git host** - GitHub, GitLab, Bitbucket, or self-hosted
- 🎯 **Focused** - Does one thing and does it well

## How It Works

When someone runs `go get your.domain/package`, the Go tool makes an HTTP request to your vanity server. The server responds with HTML containing special meta tags that tell Go where to find the actual repository:

```html
<meta name="go-import" content="go.gllm.dev/mypackage git https://github.com/gllm-dev/mypackage">
<meta name="go-source" content="go.gllm.dev/mypackage https://github.com/gllm-dev/mypackage https://github.com/gllm-dev/mypackage/tree/main{/dir} https://github.com/gllm-dev/mypackage/blob/main{/dir}/{file}#L{line}">
```

## Quick Start

### Using Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e VANITY_DOMAIN=go.gllm.dev \
  -e VANITY_REPOSITORY=https://github.com/gllm-dev \
  vanity-go:latest
```

### Using Docker Compose

```yaml
version: '3.8'

services:
  vanity-go:
    image: ghcr.io/gllm-dev/vanity-go:latest
    ports:
      - "8080:8080"
    environment:
      - VANITY_DOMAIN=go.gllm.dev
      - VANITY_REPOSITORY=https://github.com/gllm-dev
    restart: unless-stopped
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/gllm-dev/vanity-go.git
cd vanity-go

# Build the binary
go build -o vanity-go cmd/main.go

# Run with environment variables
VANITY_DOMAIN=go.gllm.dev VANITY_REPOSITORY=https://github.com/gllm-dev ./vanity-go
```

## Configuration

The server requires two environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `VANITY_DOMAIN` | Your vanity domain | `go.gllm.dev` |
| `VANITY_REPOSITORY` | Base repository URL | `https://github.com/gllm-dev` |
| `PORT` | Server port (optional) | `8080` (default) |

## Deployment

### Deployment on Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vanity-go
spec:
  replicas: 2
  selector:
    matchLabels:
      app: vanity-go
  template:
    metadata:
      labels:
        app: vanity-go
    spec:
      containers:
      - name: vanity-go
        image: ghcr.io/gllm-dev/vanity-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: VANITY_DOMAIN
          value: "go.gllm.dev"
        - name: VANITY_REPOSITORY
          value: "https://github.com/gllm-dev"
---
apiVersion: v1
kind: Service
metadata:
  name: vanity-go
spec:
  selector:
    app: vanity-go
  ports:
  - port: 80
    targetPort: 8080
```

### Nginx Configuration

If you're using Nginx as a reverse proxy:

```nginx
server {
    listen 80;
    server_name go.gllm.dev;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Usage Examples

Once deployed, users can import your packages using your custom domain:

```bash
# Install a package
go get go.gllm.dev/mypackage

# Import in Go code
import "go.gllm.dev/mypackage"
```

The server will automatically handle any package path under your domain:
- `go.gllm.dev/pkg1` → `https://github.com/gllm-dev/pkg1`
- `go.gllm.dev/pkg2` → `https://github.com/gllm-dev/pkg2`
- `go.gllm.dev/tools/cli` → `https://github.com/gllm-dev/tools`

## Development

### Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized development)

### Running Tests

```bash
go test ./...
```

### Building Docker Image

```bash
docker build -t vanity-go:latest .
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the simplicity of Go's vanity import path system
- Built with love for the Go community

## Support

- 🐛 [Report bugs](https://github.com/gllm-dev/vanity-go/issues)
- 💡 [Request features](https://github.com/gllm-dev/vanity-go/issues)
- 📖 [Read the docs](https://github.com/gllm-dev/vanity-go/wiki)

---

Made with ❤️ by the Go community