package catalog

import (
	"context"
	"database/sql"
	"eneba-marketplace/backend/internal/httpjson"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func integerParam(r *http.Request, name string, def int) (int, error) {
	s := r.URL.Query().Get(name)
	if s == "" {
		return def, nil
	}
	return strconv.Atoi(s)
}

type Lister interface {
	List(context.Context, Filters) (ProductsResponse, error)
}

func ListHandler(repo Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, e1 := integerParam(r, "page", 1)
		limit, e2 := integerParam(r, "limit", 12)
		min, e3 := integerParam(r, "min_price", 0)
		max, e4 := integerParam(r, "max_price", 0)
		f := Filters{Search: strings.TrimSpace(r.URL.Query().Get("q")), Platform: r.URL.Query().Get("platform"), Region: r.URL.Query().Get("region"), Sort: r.URL.Query().Get("sort"), MinPrice: int64(min), MaxPrice: int64(max), Page: page, Limit: limit}
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || f.Page < 1 || f.Limit < 1 || f.Limit > 24 || Validate(f) != nil {
			httpjson.Error(w, 400, "Invalid catalog filters")
			return
		}
		ctx, cancel := contextTimeout(r, 5*time.Second)
		defer cancel()
		result, err := repo.List(ctx, f)
		if err != nil {
			if err.Error() == "invalid sort" {
				httpjson.Error(w, 400, "Invalid sort")
				return
			}
			log.Printf("catalog list: %v", err)
			httpjson.Error(w, 500, "Products unavailable")
			return
		}
		httpjson.Write(w, 200, result)
	}
}
func DetailHandler(repo SQLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if len(slug) > 160 || slug == "" {
			httpjson.Error(w, 400, "Invalid product")
			return
		}
		ctx, cancel := contextTimeout(r, 5*time.Second)
		defer cancel()
		item, err := repo.Get(ctx, slug)
		if errors.Is(err, sql.ErrNoRows) {
			httpjson.Error(w, 404, "Product not found")
			return
		}
		if err != nil {
			log.Printf("product detail: %v", err)
			httpjson.Error(w, 500, "Product unavailable")
			return
		}
		httpjson.Write(w, 200, item)
	}
}
func SuggestHandler(repo SQLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if len(q) < 2 || len(q) > 100 {
			httpjson.Write(w, 200, map[string]any{"items": []Product{}})
			return
		}
		ctx, cancel := contextTimeout(r, 4*time.Second)
		defer cancel()
		out, err := repo.List(ctx, Filters{Search: q, Limit: 5, Sort: "popular"})
		if err != nil {
			httpjson.Error(w, 500, "Search unavailable")
			return
		}
		httpjson.Write(w, 200, map[string]any{"items": out.Items})
	}
}
