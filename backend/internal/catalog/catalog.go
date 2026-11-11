package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Product struct {
	ID               int64  `json:"id"`
	Slug             string `json:"slug"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Platform         string `json:"platform"`
	Region           string `json:"region"`
	Genre            string `json:"genre"`
	CoverImageURL    string `json:"cover_image_url"`
	SellerCount      int    `json:"seller_count"`
	LowestPriceCents *int64 `json:"lowest_price_cents"`
	StockCount       int    `json:"stock_count"`
	SoldCount        int    `json:"sold_count"`
}
type Offer struct {
	ID         int64  `json:"id"`
	SellerID   int64  `json:"seller_id"`
	SellerName string `json:"seller_name"`
	PriceCents int64  `json:"price_cents"`
	Stock      int    `json:"stock"`
}
type Detail struct {
	Product
	Offers []Offer `json:"offers"`
}
type Filters struct {
	Search, Platform, Region, Sort string
	MinPrice, MaxPrice             int64
	Page, Limit                    int
}
type ProductsResponse struct {
	Items    []Product `json:"items"`
	Count    int       `json:"count"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Pages    int       `json:"pages"`
}
type SQLRepository struct{ DB *sql.DB }

const baseQuery = `SELECT p.id,p.slug,p.title,p.description,p.platform,p.region,p.genre,COALESCE(p.cover_image_url,'') AS cover_image_url,
(SELECT COUNT(DISTINCT o.id) FROM offers o JOIN digital_keys k ON k.offer_id=o.id WHERE o.product_id=p.id AND o.active AND k.order_item_id IS NULL)::integer AS seller_count,
(SELECT MIN(o.price_cents) FROM offers o JOIN digital_keys k ON k.offer_id=o.id WHERE o.product_id=p.id AND o.active AND k.order_item_id IS NULL) AS lowest_price_cents,
(SELECT COUNT(*) FROM offers o JOIN digital_keys k ON k.offer_id=o.id WHERE o.product_id=p.id AND o.active AND k.order_item_id IS NULL)::integer AS stock_count,
(SELECT COUNT(*) FROM order_items oi JOIN offers o ON o.id=oi.offer_id JOIN orders ord ON ord.id=oi.order_id WHERE o.product_id=p.id AND ord.status='paid')::integer AS sold_count
FROM products p`

func scanProduct(scan func(...any) error) (Product, error) {
	var p Product
	var price sql.NullInt64
	err := scan(&p.ID, &p.Slug, &p.Title, &p.Description, &p.Platform, &p.Region, &p.Genre, &p.CoverImageURL, &p.SellerCount, &price, &p.StockCount, &p.SoldCount)
	if price.Valid {
		v := price.Int64
		p.LowestPriceCents = &v
	}
	return p, err
}
func (r SQLRepository) List(ctx context.Context, f Filters) (ProductsResponse, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 24 {
		f.Limit = 12
	}
	if f.Sort == "" {
		f.Sort = "popular"
	}
	order := map[string]string{
		"popular": "v.sold_count DESC, v.id ASC", "price_asc": "v.lowest_price_cents ASC NULLS LAST, v.id ASC",
		"price_desc": "v.lowest_price_cents DESC NULLS LAST, v.id ASC", "title": "v.title ASC, v.id ASC", "newest": "v.id DESC",
	}[f.Sort]
	if order == "" {
		return ProductsResponse{}, errors.New("invalid sort")
	}
	conditions := ` WHERE ($1='' OR LOWER(v.title) LIKE '%' || LOWER($1) || '%' OR similarity(LOWER(v.title),LOWER($1)) > 0.2)
 AND ($2='' OR v.platform=$2) AND ($3='' OR v.region=$3)
 AND ($4=0 OR v.lowest_price_cents >= $4) AND ($5=0 OR v.lowest_price_cents <= $5)`
	args := []any{f.Search, f.Platform, f.Region, f.MinPrice, f.MaxPrice}
	cte := `WITH v AS (` + baseQuery + `) `
	var total int
	if err := r.DB.QueryRowContext(ctx, cte+`SELECT COUNT(*) FROM v`+conditions, args...).Scan(&total); err != nil {
		return ProductsResponse{}, fmt.Errorf("count products: %w", err)
	}
	query := cte + `SELECT v.id,v.slug,v.title,v.description,v.platform,v.region,v.genre,v.cover_image_url,v.seller_count,v.lowest_price_cents,v.stock_count,v.sold_count FROM v` + conditions + ` ORDER BY ` + order + ` LIMIT $6 OFFSET $7`
	rows, err := r.DB.QueryContext(ctx, query, append(args, f.Limit, (f.Page-1)*f.Limit)...)
	if err != nil {
		return ProductsResponse{}, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()
	items := make([]Product, 0)
	for rows.Next() {
		p, e := scanProduct(rows.Scan)
		if e != nil {
			return ProductsResponse{}, e
		}
		items = append(items, p)
	}
	if err = rows.Err(); err != nil {
		return ProductsResponse{}, err
	}
	return ProductsResponse{Items: items, Count: len(items), Total: total, Page: f.Page, PageSize: f.Limit, Pages: (total + f.Limit - 1) / f.Limit}, nil
}
func (r SQLRepository) Get(ctx context.Context, slug string) (Detail, error) {
	row := r.DB.QueryRowContext(ctx, baseQuery+` WHERE p.slug=$1`, slug)
	p, err := scanProduct(row.Scan)
	if err != nil {
		return Detail{}, err
	}
	rows, err := r.DB.QueryContext(ctx, `SELECT o.id,s.id,s.name,o.price_cents,COUNT(k.id)::int
 FROM offers o JOIN sellers s ON s.id=o.seller_id
 LEFT JOIN digital_keys k ON k.offer_id=o.id AND k.order_item_id IS NULL
 WHERE o.product_id=$1 AND o.active GROUP BY o.id,s.id,s.name ORDER BY o.price_cents,o.id`, p.ID)
	if err != nil {
		return Detail{}, err
	}
	defer rows.Close()
	offers := make([]Offer, 0)
	for rows.Next() {
		var o Offer
		if err := rows.Scan(&o.ID, &o.SellerID, &o.SellerName, &o.PriceCents, &o.Stock); err != nil {
			return Detail{}, err
		}
		offers = append(offers, o)
	}
	if err = rows.Err(); err != nil {
		return Detail{}, err
	}
	return Detail{Product: p, Offers: offers}, nil
}
func Validate(f Filters) error {
	if len([]rune(f.Search)) > 100 || len(f.Platform) > 80 || len(f.Region) > 80 {
		return errors.New("search/filter too long")
	}
	if f.MinPrice < 0 || f.MaxPrice < 0 || (f.MaxPrice > 0 && f.MinPrice > f.MaxPrice) {
		return errors.New("invalid price range")
	}
	if f.Page > 10000 {
		return errors.New("page too large")
	}
	if strings.ContainsRune(f.Search, 0) {
		return errors.New("invalid search")
	}
	return nil
}
