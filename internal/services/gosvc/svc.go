package gosvc

import (
	"context"
	"strings"
)

// Service handles the generation of vanity import HTML responses.
// It contains the domain and repository information needed to construct
// the proper meta tags for Go's import path resolution.
type Service struct {
	// domain is the vanity domain (e.g., "go.gllm.dev")
	domain string
	// repository is the base repository URL (e.g., "https://github.com/gllm-dev")
	repository string
}

// New creates a new Service instance with the given domain and repository base URL.
// The domain should be the vanity import domain without protocol (e.g., "go.gllm.dev").
// The repository should be the base URL where modules are hosted (e.g., "https://github.com/gllm-dev").
func New(domain, repository string) *Service {
	return &Service{
		domain:     domain,
		repository: repository,
	}
}

// template defines the HTML template returned for vanity import requests.
// It includes:
// - go-import meta tag: tells go get where to find the repository
// - go-source meta tag: provides source code browsing information for godoc.org
// The placeholders {{.domain}}, {{.repository}}, and {{.module}} are replaced
// with actual values when generating the response.
const template = `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="{{.domain}} git {{.repository}}">
<meta name="go-source" content="{{.domain}} {{.repository}} {{.repository}}/tree/main{/dir} {{.repository}}/blob/main{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here; <a href="https://pkg.go.dev/{{.domain}}">see the package on pkg.go.dev</a>.
</body>
</html>`

// Vanity generates the HTML response for a given module path.
// It takes the module name and returns an HTML string with the appropriate
// go-import and go-source meta tags for Go's import path resolution.
//
// The generated HTML allows `go get` to resolve custom import paths like
// "go.gllm.dev/vanity-go" to the actual repository location.
//
// Example:
//
//	For domain="go.gllm.dev", repository="https://github.com/gllm-dev", and module="vanity-go",
//	it generates meta tags that redirect "go.gllm.dev/vanity-go" to "https://github.com/gllm-dev/vanity-go".
func (s *Service) Vanity(ctx context.Context, module string) string {
	domain := s.domain
	repository := s.repository

	var fullDomain, fullRepository string
	if module == "" {
		fullDomain = domain
		fullRepository = repository
	} else {
		fullDomain = domain + "/" + module
		fullRepository = repository + "/" + module
	}

	parsedTemplate := strings.ReplaceAll(template, "{{.domain}}", fullDomain)
	return strings.ReplaceAll(parsedTemplate, "{{.repository}}", fullRepository)
}
