package homepage

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"eneba-marketplace/backend/internal/catalog"
	"eneba-marketplace/backend/internal/httpjson"
)

type Banner struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
	Eyebrow  string `json:"eyebrow"`
	Action   string `json:"action"`
	Href     string `json:"href"`
}
type Collection struct {
	Heading  string            `json:"heading"`
	Subtitle string            `json:"subtitle"`
	Items    []catalog.Product `json:"items"`
}

func BannerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var b Banner
		err := db.QueryRowContext(ctx, "SELECT heading,subtitle FROM homepage_sections WHERE section_key='banner'").Scan(&b.Heading, &b.Subtitle)
		if err != nil {
			log.Printf("banner: %v", err)
			httpjson.Error(w, 500, "Banner unavailable")
			return
		}
		b.Eyebrow = "YOUR NEXT ADVENTURE STARTS HERE"
		b.Action = "Explore games"
		b.Href = "/store"
		httpjson.Write(w, 200, b)
	}
}
func CollectionHandler(db *sql.DB, kind string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var c Collection
		err := db.QueryRowContext(ctx, `SELECT heading,subtitle FROM homepage_sections WHERE section_key=$1`, kind).Scan(&c.Heading, &c.Subtitle)
		if err != nil {
			httpjson.Error(w, 500, "Collection unavailable")
			return
		}
		sort := "popular"
		if kind == "deals" {
			sort = "price_asc"
		}
		data, err := (catalog.SQLRepository{DB: db}).List(ctx, catalog.Filters{Sort: sort, Limit: 8})
		if err != nil {
			log.Printf("collection: %v", err)
			httpjson.Error(w, 500, "Collection unavailable")
			return
		}
		c.Items = data.Items
		httpjson.Write(w, 200, c)
	}
}
