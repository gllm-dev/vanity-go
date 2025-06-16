# vanity-go Examples

This directory contains example configurations for deploying vanity-go in various environments.

## Files

### `.env.example`
Example environment variables file showing the required configuration options. Copy this to `.env` and modify the values to match your setup.

### `docker-compose.yml`
Docker Compose configuration for running vanity-go with Docker. Includes:
- Basic service configuration
- Health checks
- Optional resource limits
- Example Traefik labels for reverse proxy

### `kubernetes.yaml`
Complete Kubernetes manifests including:
- Deployment with security best practices
- Service for internal cluster access
- Ingress for external HTTPS access
- ConfigMap for configuration management
- HorizontalPodAutoscaler for automatic scaling
- PodDisruptionBudget for high availability

### `systemd.service`
Systemd service file for running vanity-go on Linux systems. Includes:
- Security hardening options
- Resource limits
- Automatic restart configuration
- Proper logging setup

### `nginx.conf`
Nginx reverse proxy configuration with:
- HTTP to HTTPS redirect
- Modern SSL/TLS settings
- Security headers
- Proxy configuration optimized for Go imports
- Optional load balancing setup

## Quick Start Examples

### Local Development
```bash
# Using environment variables
VANITY_DOMAIN=localhost:8080 VANITY_REPOSITORY=https://github.com/gllm-dev ./vanity-go

# Using .env file
cp .env.example .env
# Edit .env with your values
./vanity-go
```

### Docker
```bash
# Single container
docker run -d \
  -p 8080:8080 \
  -e VANITY_DOMAIN=go.gllm.dev \
  -e VANITY_REPOSITORY=https://github.com/gllm-dev \
  vanity-go:latest

# Using docker-compose
docker-compose up -d
```

### Kubernetes
```bash
# Update the YAML with your values
kubectl apply -f kubernetes.yaml

# Check deployment status
kubectl get pods -l app=vanity-go
kubectl get ingress vanity-go
```

### Systemd (Linux)
```bash
# Copy binary
sudo cp vanity-go /opt/vanity-go/

# Create user
sudo useradd -r -s /bin/false vanity-go

# Install service
sudo cp systemd.service /etc/systemd/system/vanity-go.service
sudo systemctl daemon-reload
sudo systemctl enable vanity-go
sudo systemctl start vanity-go

# Check status
sudo systemctl status vanity-go
sudo journalctl -u vanity-go -f
```

### Nginx Reverse Proxy
```bash
# Copy configuration
sudo cp nginx.conf /etc/nginx/sites-available/go.gllm.dev
sudo ln -s /etc/nginx/sites-available/go.gllm.dev /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

## Production Considerations

### Security
- Always use HTTPS in production
- Run as non-root user
- Use read-only root filesystem where possible
- Implement rate limiting at the reverse proxy level

### Performance
- Enable caching at the reverse proxy (5-10 minutes is reasonable)
- Use horizontal scaling for high traffic
- Monitor resource usage and adjust limits accordingly

### Monitoring
- Set up health checks
- Monitor HTTP response times
- Track 4xx/5xx errors
- Set up alerts for service downtime

### Backup and Recovery
- vanity-go is stateless, so no data backup needed
- Keep configuration in version control
- Document your domain and repository mappings

## Common Patterns

### Multiple Repositories
If you have packages in multiple repositories, you can run multiple instances:

```yaml
# docker-compose.yml
services:
  vanity-github:
    image: vanity-go:latest
    environment:
      - VANITY_DOMAIN=go.gllm.dev/github
      - VANITY_REPOSITORY=https://github.com/gllm-dev
  
  vanity-gitlab:
    image: vanity-go:latest
    environment:
      - VANITY_DOMAIN=go.gllm.dev/gitlab
      - VANITY_REPOSITORY=https://gitlab.com/gllm-dev
```

### Subdomain per Organization
```nginx
# Nginx: Route by subdomain
server {
    server_name go.team1.example.com;
    location / {
        proxy_pass http://team1-vanity:8080;
    }
}

server {
    server_name go.team2.example.com;
    location / {
        proxy_pass http://team2-vanity:8080;
    }
}
```

## Troubleshooting

### Test the vanity import
```bash
# Check if the server is responding correctly
curl -H "User-Agent: go-get" https://go.gllm.dev/mypackage?go-get=1

# Test actual go get
go get -v go.gllm.dev/mypackage
```

### Common Issues

1. **404 errors**: Check that the repository exists and is accessible
2. **SSL errors**: Ensure certificates are valid and properly configured
3. **Timeout errors**: Check network connectivity and firewall rules
4. **Import not working**: Verify the HTML meta tags are correct

For more help, see the main README or open an issue on GitHub.