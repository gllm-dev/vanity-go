package gosvc

import (
	"context"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		repository string
	}{
		{
			name:       "basic initialization",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
		},
		{
			name:       "domain with subdomain",
			domain:     "pkg.company.com",
			repository: "https://gitlab.com/company",
		},
		{
			name:       "repository with trailing slash",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.domain, tt.repository)
			if svc == nil {
				t.Fatal("expected non-nil service")
			}
			if svc.domain != tt.domain {
				t.Errorf("domain = %v, want %v", svc.domain, tt.domain)
			}
			if svc.repository != tt.repository {
				t.Errorf("repository = %v, want %v", svc.repository, tt.repository)
			}
		})
	}
}

func TestService_Vanity(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		repository string
		pkg        string
		wantChecks []string
	}{
		{
			name:       "root package",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "",
			wantChecks: []string{
				`<meta name="go-import" content="go.gllm.dev git https://github.com/gllm-dev">`,
				`<meta name="go-source" content="go.gllm.dev https://github.com/gllm-dev https://github.com/gllm-dev/tree/main{/dir} https://github.com/gllm-dev/blob/main{/dir}/{file}#L{line}">`,
				`<!DOCTYPE html>`,
				`<html>`,
				`</html>`,
			},
		},
		{
			name:       "sub package",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "tools",
			wantChecks: []string{
				`<meta name="go-import" content="go.gllm.dev/tools git https://github.com/gllm-dev/tools">`,
				`<meta name="go-source" content="go.gllm.dev/tools https://github.com/gllm-dev/tools https://github.com/gllm-dev/tools/tree/main{/dir} https://github.com/gllm-dev/tools/blob/main{/dir}/{file}#L{line}">`,
			},
		},
		{
			name:       "nested package",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "cmd/tool",
			wantChecks: []string{
				`<meta name="go-import" content="go.gllm.dev/cmd/tool git https://github.com/gllm-dev/cmd/tool">`,
				`<meta name="go-source" content="go.gllm.dev/cmd/tool https://github.com/gllm-dev/cmd/tool https://github.com/gllm-dev/cmd/tool/tree/main{/dir} https://github.com/gllm-dev/cmd/tool/blob/main{/dir}/{file}#L{line}">`,
			},
		},
		{
			name:       "repository with trailing slash",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev/",
			pkg:        "pkg",
			wantChecks: []string{
				`<meta name="go-import" content="go.gllm.dev/pkg git https://github.com/gllm-dev//pkg">`,
			},
		},
		{
			name:       "GitLab repository",
			domain:     "go.company.com",
			repository: "https://gitlab.com/company",
			pkg:        "internal/api",
			wantChecks: []string{
				`<meta name="go-import" content="go.company.com/internal/api git https://gitlab.com/company/internal/api">`,
				`<meta name="go-source" content="go.company.com/internal/api https://gitlab.com/company/internal/api https://gitlab.com/company/internal/api/tree/main{/dir} https://gitlab.com/company/internal/api/blob/main{/dir}/{file}#L{line}">`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.domain, tt.repository)
			got := svc.Vanity(context.Background(), tt.pkg)

			for _, check := range tt.wantChecks {
				if !strings.Contains(got, check) {
					t.Errorf("Vanity() missing expected content:\nwant substring: %s\ngot: %s", check, got)
				}
			}

			// Check that it's valid HTML structure
			if !strings.HasPrefix(got, "<!DOCTYPE html>") {
				t.Error("Vanity() should start with <!DOCTYPE html>")
			}
			if !strings.Contains(got, "<html>") {
				t.Error("Vanity() should contain <html> tag")
			}
			if !strings.Contains(got, "</html>") {
				t.Error("Vanity() should contain </html> tag")
			}
			if !strings.Contains(got, "<head>") {
				t.Error("Vanity() should contain <head> tag")
			}
			if !strings.Contains(got, "</head>") {
				t.Error("Vanity() should contain </head> tag")
			}
			if !strings.Contains(got, "<body>") {
				t.Error("Vanity() should contain <body> tag")
			}
			if !strings.Contains(got, "</body>") {
				t.Error("Vanity() should contain </body> tag")
			}
		})
	}
}

func TestService_Vanity_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		repository string
		pkg        string
		wantError  bool
	}{
		{
			name:       "package with hyphen",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "my-package",
			wantError:  false,
		},
		{
			name:       "package with underscore",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "my_package",
			wantError:  false,
		},
		{
			name:       "package with dot",
			domain:     "go.gllm.dev",
			repository: "https://github.com/gllm-dev",
			pkg:        "v2.0",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.domain, tt.repository)
			got := svc.Vanity(context.Background(), tt.pkg)

			if tt.wantError {
				// Currently the service doesn't return errors, but this is here for future enhancement
				t.Skip("Error handling not yet implemented")
			}

			expectedImport := tt.domain + "/" + tt.pkg
			if !strings.Contains(got, expectedImport) {
				t.Errorf("Vanity() should contain import path %s", expectedImport)
			}
		})
	}
}

func BenchmarkService_Vanity(b *testing.B) {
	svc := New("go.gllm.dev", "https://github.com/gllm-dev")
	
	b.Run("root_package", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = svc.Vanity(context.Background(), "")
		}
	})

	b.Run("sub_package", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = svc.Vanity(context.Background(), "pkg/subpkg")
		}
	})

	b.Run("deep_package", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = svc.Vanity(context.Background(), "pkg/sub/deep/nested/package")
		}
	})
}