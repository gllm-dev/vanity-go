package gohdl

import (
	"go.gllm.dev/vanity-go/internal/services/gosvc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	svc := gosvc.New("go.gllm.dev", "https://github.com/gllm-dev")
	h := New(svc)

	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.service == nil {
		t.Fatal("expected non-nil service in handler")
	}
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		requestPath    string
		queryParams    string
		wantStatusCode int
		wantContains   []string
		wantHeader     map[string]string
	}{
		{
			name:           "root path",
			requestPath:    "/",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev git https://github.com/gllm-dev">`,
				`<!DOCTYPE html>`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "package path",
			requestPath:    "/mypackage",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/mypackage git https://github.com/gllm-dev/mypackage">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "nested package path",
			requestPath:    "/cmd/tool/cli",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/cmd/tool/cli git https://github.com/gllm-dev/cmd/tool/cli">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "path with go-get query",
			requestPath:    "/mypackage",
			queryParams:    "?go-get=1",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/mypackage git https://github.com/gllm-dev/mypackage">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "path with trailing slash",
			requestPath:    "/mypackage/",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/mypackage/ git https://github.com/gllm-dev/mypackage/">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "package with hyphen",
			requestPath:    "/my-package",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/my-package git https://github.com/gllm-dev/my-package">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "package with underscore",
			requestPath:    "/my_package",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/my_package git https://github.com/gllm-dev/my_package">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
		{
			name:           "version path",
			requestPath:    "/v2",
			queryParams:    "",
			wantStatusCode: http.StatusOK,
			wantContains: []string{
				`<meta name="go-import" content="go.gllm.dev/v2 git https://github.com/gllm-dev/v2">`,
			},
			wantHeader: map[string]string{
				"Content-Type": "text/html; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service and handler
			svc := gosvc.New("go.gllm.dev", "https://github.com/gllm-dev")
			h := New(svc)

			// Create request
			req, err := http.NewRequest("GET", tt.requestPath+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			h.Handle(rr, req)

			// Check status code
			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

			// Check response body
			body := rr.Body.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(body, want) {
					t.Errorf("handler returned body missing expected content:\nwant: %s\ngot: %s",
						want, body)
				}
			}

			// Check headers
			for header, want := range tt.wantHeader {
				if got := rr.Header().Get(header); got != want {
					t.Errorf("handler returned wrong header %s: got %v want %v",
						header, got, want)
				}
			}
		})
	}
}

func TestHandler_Handle_DifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			svc := gosvc.New("go.gllm.dev", "https://github.com/gllm-dev")
			h := New(svc)

			req, err := http.NewRequest(method, "/package", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.Handle(rr, req)

			// All methods should return 200 OK with the vanity import HTML
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code for %s: got %v want %v",
					method, status, http.StatusOK)
			}

			// Should contain the import meta tag
			body := rr.Body.String()
			if !strings.Contains(body, "go-import") {
				t.Errorf("handler should return go-import meta tag for %s method", method)
			}
		})
	}
}

func TestHandler_Handle_SpecialPaths(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		wantImportPath string
	}{
		{
			name:           "double slash - root path",
			path:           "//package",
			wantImportPath: "go.gllm.dev git https://github.com/gllm-dev", // Becomes root path
		},
		{
			name:           "multiple trailing slashes",
			path:           "/package///",
			wantImportPath: "go.gllm.dev/package/// git https://github.com/gllm-dev/package///",
		},
		{
			name:           "encoded slash gets decoded by net/url",
			path:           "/package%2Fsub",
			wantImportPath: "go.gllm.dev/package/sub git https://github.com/gllm-dev/package/sub", // URL decoding happens automatically
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := gosvc.New("go.gllm.dev", "https://github.com/gllm-dev")
			h := New(svc)

			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.Handle(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			// Check that the response contains expected import path
			body := rr.Body.String()
			if !strings.Contains(body, tt.wantImportPath) {
				t.Errorf("handler should contain import path %s, got: %s", tt.wantImportPath, body)
			}
		})
	}
}

func BenchmarkHandler_Handle(b *testing.B) {
	svc := gosvc.New("go.gllm.dev", "https://github.com/gllm-dev")
	h := New(svc)

	paths := []string{
		"/",
		"/package",
		"/deep/nested/package/path",
	}

	for _, path := range paths {
		b.Run(strings.ReplaceAll(path, "/", "_"), func(b *testing.B) {
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rr := httptest.NewRecorder()
				h.Handle(rr, req)
			}
		})
	}
}
