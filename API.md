# vanity-go API Documentation

vanity-go provides a simple HTTP API that responds to Go tool requests for package imports.

## Base URL

The base URL is your configured domain. For example:
```
https://go.gllm.dev
```

## Endpoints

### GET /{package-path}

Returns HTML with meta tags for Go import path resolution.

#### Parameters

- **package-path** (path parameter): The Go package path being requested
  - Can be empty for root packages
  - Can contain slashes for nested packages (e.g., `cmd/tool`)
  - Examples: `/`, `/mypackage`, `/tools/cli`, `/v2/api`

- **go-get** (query parameter, optional): Set to `1` by the Go tool
  - The server responds identically whether this parameter is present or not
  - Example: `/mypackage?go-get=1`

#### Response

**Status Code:** 200 OK

**Content-Type:** text/html; charset=utf-8

**Body:** HTML document containing go-import and go-source meta tags

#### Example Request

```bash
GET /vanity-go?go-get=1
Host: go.gllm.dev
```

#### Example Response

```html
<!DOCTYPE html>
<html>
<head>
    <meta name="go-import" content="go.gllm.dev/vanity-go git https://github.com/gllm-dev/vanity-go">
    <meta name="go-source" content="go.gllm.dev/vanity-go https://github.com/gllm-dev/vanity-go https://github.com/gllm-dev/vanity-go/tree/main{/dir} https://github.com/gllm-dev/vanity-go/blob/main{/dir}/{file}#L{line}">
</head>
<body>
    go get go.gllm.dev/vanity-go
</body>
</html>
```

### GET /healthz

Health check endpoint for monitoring.

#### Response

**Status Code:** 200 OK

**Body:** `OK`

#### Example Request

```bash
GET /healthz
Host: go.gllm.dev
```

## Meta Tags

The HTML response includes two important meta tags:

### go-import

Tells the Go tool where to find the repository.

Format:
```html
<meta name="go-import" content="{import-path} {vcs} {repo-url}">
```

- **import-path**: The full import path (domain + package path)
- **vcs**: Version control system (always `git` in vanity-go)
- **repo-url**: The actual repository URL

### go-source

Provides information for godoc.org and pkg.go.dev to display source code.

Format:
```html
<meta name="go-source" content="{import-path} {home} {directory} {file}">
```

- **import-path**: The full import path
- **home**: Repository home page URL
- **directory**: URL template for directory browsing
- **file**: URL template for file viewing with line numbers

## How Go Uses This API

When you run:
```bash
go get go.gllm.dev/vanity-go
```

1. Go makes an HTTP GET request to `https://go.gllm.dev/vanity-go?go-get=1`
2. The server responds with HTML containing the meta tags
3. Go parses the `go-import` meta tag to find the repository URL
4. Go clones/fetches from the actual repository (e.g., GitHub)

## URL Path Mapping

The server maps URL paths to repository paths:

| Request Path | Import Path | Repository URL |
|-------------|-------------|----------------|
| `/` | `go.gllm.dev` | `https://github.com/gllm-dev` |
| `/pkg` | `go.gllm.dev/pkg` | `https://github.com/gllm-dev/pkg` |
| `/tools/cli` | `go.gllm.dev/tools/cli` | `https://github.com/gllm-dev/tools/cli` |
| `/v2` | `go.gllm.dev/v2` | `https://github.com/gllm-dev/v2` |

## Error Handling

The server always returns 200 OK with valid HTML, even for non-existent packages. This is by design:

- The Go tool will attempt to fetch from the provided repository URL
- If the repository doesn't exist, the Go tool will report the error
- This allows dynamic package creation without server updates

## Caching

Responses can be cached safely:

- The meta tags don't change frequently
- Recommended cache time: 5-10 minutes
- Use `Cache-Control: public, max-age=300` header

## Rate Limiting

The server does not implement rate limiting. It's recommended to:

- Implement rate limiting at the reverse proxy level
- Use CDN for caching and DDoS protection
- Monitor for abuse patterns

## Examples

### Using curl

```bash
# Basic request
curl https://go.gllm.dev/vanity-go

# With go-get parameter (same response)
curl https://go.gllm.dev/vanity-go?go-get=1

# Check health
curl https://go.gllm.dev/healthz
```

### Using Go

```bash
# Get a package
go get go.gllm.dev/vanity-go

# Get a specific version
go get go.gllm.dev/vanity-go@v1.2.3

# Get latest from main branch
go get go.gllm.dev/vanity-go@main
```

### Debugging

To see what Go is doing:
```bash
go get -v -x go.gllm.dev/vanity-go
```

## Security Considerations

1. **HTTPS Only**: Always use HTTPS in production to prevent MITM attacks
2. **No Authentication**: The server provides public information only
3. **Input Validation**: Package paths are passed directly to templates
4. **No State**: Server is stateless, reducing attack surface

## Performance

- Responses are generated dynamically but are very lightweight
- No database queries or external API calls
- Typical response time: < 5ms
- Can handle thousands of requests per second on modest hardware