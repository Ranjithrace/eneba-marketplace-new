package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewHandler(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(dir, "index.html"),
		[]byte("REACT PAGE"), 0600,
	); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(filepath.Join(dir, "assets"), 0700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		filepath.Join(dir, "assets", "app.js"),
		[]byte("APP SCRIPT"), 0600,
	); err != nil {
		t.Fatal(err)
	}

	api := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("API RESPONSE"))
	})

	site := NewHandler(api, dir)

	tests := []struct {
		name, method, url, contains string
		status                      int
	}{
		{"homepage", "GET", "/", "REACT PAGE", 200},
		{"store direct link", "GET", "/store", "REACT PAGE", 200},
		{"orders direct link", "GET", "/orders", "REACT PAGE", 200},
		{"compiled asset", "GET", "/assets/app.js", "APP SCRIPT", 200},
		{"missing asset", "GET", "/assets/missing.js", "", 404},
		{"missing file", "GET", "/missing.js", "", 404},
		{"API forwarding", "GET", "/api/health", "API RESPONSE", 200},
		{"reject page POST", "POST", "/store", "", 405},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(tc.method, tc.url, nil)
			site.ServeHTTP(response, request)

			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d",
					response.Code, tc.status)
			}

			if !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("body = %q, expected %q",
					response.Body.String(), tc.contains)
			}
		})
	}
}
