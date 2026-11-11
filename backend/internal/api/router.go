package api

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"eneba-marketplace/backend/internal/auth"
	"eneba-marketplace/backend/internal/catalog"
	"eneba-marketplace/backend/internal/health"
	"eneba-marketplace/backend/internal/homepage"
	"eneba-marketplace/backend/internal/httpjson"
	"eneba-marketplace/backend/internal/orders"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	catalogRepo := catalog.SQLRepository{DB: db}
	secure := os.Getenv("COOKIE_SECURE") == "true"
	account := auth.Service{DB: db, SecureCookies: secure}
	purchasing := orders.Service{DB: db}
	mux.HandleFunc("GET /api/health", health.Handler(db))
	mux.HandleFunc("GET /api/home/banner", homepage.BannerHandler(db))
	mux.HandleFunc("GET /api/home/bestsellers", homepage.CollectionHandler(db, "bestsellers"))
	mux.HandleFunc("GET /api/home/deals", homepage.CollectionHandler(db, "deals"))
	mux.HandleFunc("GET /api/products", catalog.ListHandler(catalogRepo))
	mux.HandleFunc("GET /api/products/{slug}", catalog.DetailHandler(catalogRepo))
	mux.HandleFunc("GET /api/search/suggest", catalog.SuggestHandler(catalogRepo))
	mux.HandleFunc("POST /api/auth/register", account.Register)
	mux.HandleFunc("POST /api/auth/login", account.Login)
	mux.HandleFunc("POST /api/auth/logout", account.Logout)
	mux.HandleFunc("GET /api/auth/me", account.Me)
	mux.HandleFunc("POST /api/checkout", account.Required(purchasing.Checkout))
	mux.HandleFunc("GET /api/orders", account.Required(purchasing.History))
	mux.HandleFunc("GET /api/orders/{id}", account.Required(purchasing.Detail))
	return security(mux)
}
func security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			origin := r.Header.Get("Origin")
			allowed := strings.TrimRight(os.Getenv("PUBLIC_ORIGIN"), "/")
			localOK := allowed == "" && (origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173")
			if origin != "" && !localOK && origin != allowed && origin != "http://"+r.Host && origin != "https://"+r.Host {
				httpjson.Error(w, http.StatusForbidden, "Request origin not allowed")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
