# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.0.1] - 2025-06-16
### Added
- Initial release of vanity-go
- Basic vanity import server functionality
- Support for custom domain and repository configuration
- HTML meta tags generation for go-import and go-source
- Docker support with multi-stage build
- Health check endpoint at `/health`
- Configurable port via PORT environment variable
- Graceful shutdown support
- Request timeouts for better resource management
- Comprehensive documentation including:
  - README with quick start guide
  - API documentation
  - Contributing guidelines
  - Example configurations for various deployment scenarios
- Unit tests for core functionality
- GitHub Actions workflows for CI/CD
- Makefile for common development tasks