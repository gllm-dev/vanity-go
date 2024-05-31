package gosvc

import (
	"context"
	"strings"
)

type Service struct {
	domain     string
	repository string
}

func New(domain, repository string) *Service {
	return &Service{
		domain:     domain,
		repository: repository,
	}
}

const template = `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="{{.domain}}/{{.module}} git {{.repository}}/{{.module}}">
<meta name="go-source" content="{{.domain}}/{{.module}} {{.repository}}/{{.module}} {{.repository}}/{{.module}}/tree/main{/dir} {{.repository}}/{{.module}}/blob/main{/dir}/{file}#L{line}">
</head>
<body>
Nothing to see here; <a href="https://pkg.go.dev/{{.domain}}/{{.module}}">see the package on pkg.go.dev</a>.
</body>
</html>`

func (s *Service) Vanity(ctx context.Context, module string) string {
	parsedTemplate := strings.ReplaceAll(template, "{{.domain}}", s.domain)
	parsedTemplate = strings.ReplaceAll(parsedTemplate, "{{.repository}}", s.repository)
	return strings.ReplaceAll(parsedTemplate, "{{.module}}", module)
}
