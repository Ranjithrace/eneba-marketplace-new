package web

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NewHandler serves the API and compiled React app on one origin.
func NewHandler(api http.Handler, distDir string) http.Handler {
	files := http.FileServer(http.Dir(distDir))
	index := filepath.Join(distDir, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		clean := path.Clean("/" + r.URL.Path)
		relative := strings.TrimPrefix(clean, "/")

		if relative != "" {
			file, err := os.Stat(
				filepath.Join(distDir, filepath.FromSlash(relative)),
			)
			if err == nil && file.Mode().IsRegular() {
				files.ServeHTTP(w, r)
				return
			}
		}

		if strings.HasPrefix(clean, "/assets/") || path.Ext(clean) != "" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, index)
	})
}
