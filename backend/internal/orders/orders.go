package orders

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"eneba-marketplace/backend/internal/auth"
	"eneba-marketplace/backend/internal/httpjson"
)

type Line struct {
	OfferID  int64 `json:"offer_id"`
	Quantity int   `json:"quantity"`
}
type CheckoutRequest struct {
	Items          []Line `json:"items"`
	PaymentMethod  string `json:"payment_method"`
	IdempotencyKey string `json:"idempotency_key"`
}
type PurchasedItem struct {
	ID         int64   `json:"id"`
	OfferID    int64   `json:"offer_id"`
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	Seller     string  `json:"seller"`
	PriceCents int64   `json:"price_cents"`
	KeyCode    *string `json:"key_code"`
}
type Order struct {
	ID         int64           `json:"id"`
	Status     string          `json:"status"`
	TotalCents int64           `json:"total_cents"`
	CreatedAt  time.Time       `json:"created_at"`
	Items      []PurchasedItem `json:"items"`
}
type Service struct{ DB *sql.DB }

var errBadRequest = errors.New("invalid checkout")

// Normalize combines repeated offers, sorts them for stable transaction lock order,
// and caps the entire demo checkout to 20 digital keys.
func Normalize(input CheckoutRequest) (CheckoutRequest, error) {
	if input.PaymentMethod != "approve" && input.PaymentMethod != "decline" {
		return CheckoutRequest{}, errBadRequest
	}
	if len(input.IdempotencyKey) < 12 || len(input.IdempotencyKey) > 100 {
		return CheckoutRequest{}, errBadRequest
	}
	if len(input.Items) < 1 || len(input.Items) > 20 {
		return CheckoutRequest{}, errBadRequest
	}
	qty := make(map[int64]int)
	total := 0
	for _, line := range input.Items {
		if line.OfferID <= 0 || line.Quantity < 1 || line.Quantity > 20 {
			return CheckoutRequest{}, errBadRequest
		}
		total += line.Quantity
		if total > 20 {
			return CheckoutRequest{}, errBadRequest
		}
		qty[line.OfferID] += line.Quantity
	}
	input.Items = make([]Line, 0, len(qty))
	for id, n := range qty {
		input.Items = append(input.Items, Line{OfferID: id, Quantity: n})
	}
	sort.Slice(input.Items, func(i, j int) bool { return input.Items[i].OfferID < input.Items[j].OfferID })
	return input, nil
}
func fingerprint(req CheckoutRequest) string {
	b, _ := json.Marshal(struct {
		Items   []Line `json:"items"`
		Payment string `json:"payment"`
	}{req.Items, req.PaymentMethod})
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:])
}
func Fingerprint(req CheckoutRequest) (string, error) {
	n, e := Normalize(req)
	if e != nil {
		return "", e
	}
	return fingerprint(n), nil
}
func (s Service) Checkout(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFrom(r.Context())
	var req CheckoutRequest
	if !httpjson.Decode(w, r, &req) {
		return
	}
	normalized, err := Normalize(req)
	if err != nil {
		httpjson.Error(w, 400, "Provide valid offers, quantity (max 20), payment choice and an idempotency key")
		return
	}
	req = normalized
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		httpjson.Error(w, 500, "Checkout unavailable")
		return
	}
	defer tx.Rollback()
	var id int64
	err = tx.QueryRowContext(ctx, `INSERT INTO orders(user_id,status,total_cents,idempotency_key,request_hash) VALUES($1,'pending',0,$2,$3) ON CONFLICT DO NOTHING RETURNING id`, u.ID, req.IdempotencyKey, fingerprint(req)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		var existingID int64
		var hash string
		e := tx.QueryRowContext(ctx, `SELECT id,request_hash FROM orders WHERE user_id=$1 AND idempotency_key=$2`, u.ID, req.IdempotencyKey).Scan(&existingID, &hash)
		if e != nil {
			httpjson.Error(w, 409, "Checkout in progress; retry with same key")
			return
		}
		if hash != fingerprint(req) {
			httpjson.Error(w, 409, "Idempotency key already used for a different cart")
			return
		}
		order, e := loadOrder(ctx, tx, existingID, u.ID)
		if e != nil {
			httpjson.Error(w, 500, "Could not load order")
			return
		}
		if err = tx.Commit(); err != nil {
			httpjson.Error(w, 500, "Order unavailable")
			return
		}
		httpjson.Write(w, 200, order)
		return
	}
	if err != nil {
		log.Printf("checkout insert: %v", err)
		httpjson.Error(w, 500, "Checkout unavailable")
		return
	}
	var amount int64
	prices := make(map[int64]int64, len(req.Items))
	for _, line := range req.Items {
		var price int64
		// Lock each offer in ascending ID order. All checkouts for that offer serialize.
		e := tx.QueryRowContext(ctx, `SELECT price_cents FROM offers WHERE id=$1 AND active=TRUE FOR UPDATE`, line.OfferID).Scan(&price)
		if errors.Is(e, sql.ErrNoRows) {
			httpjson.Error(w, 409, "An offer is no longer available")
			return
		}
		if e != nil {
			log.Printf("offer lock: %v", e)
			httpjson.Error(w, 500, "Checkout unavailable")
			return
		}
		var available int
		e = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM digital_keys WHERE offer_id=$1 AND order_item_id IS NULL`, line.OfferID).Scan(&available)
		if e != nil {
			httpjson.Error(w, 500, "Stock unavailable")
			return
		}
		if available < line.Quantity {
			httpjson.Error(w, 409, "Not enough digital keys in stock")
			return
		}
		amount += price * int64(line.Quantity)
		prices[line.OfferID] = price
	}
	status := "paid"
	if req.PaymentMethod == "decline" {
		status = "failed"
	}
	if _, err = tx.ExecContext(ctx, `UPDATE orders SET status=$1,total_cents=$2 WHERE id=$3`, status, amount, id); err != nil {
		httpjson.Error(w, 500, "Checkout unavailable")
		return
	}
	for _, line := range req.Items {
		for i := 0; i < line.Quantity; i++ {
			var itemID int64
			err = tx.QueryRowContext(ctx, `INSERT INTO order_items(order_id,offer_id,price_cents) VALUES($1,$2,$3) RETURNING id`, id, line.OfferID, prices[line.OfferID]).Scan(&itemID)
			if err != nil {
				httpjson.Error(w, 500, "Checkout unavailable")
				return
			}
			if status == "paid" {
				var keyID int64
				err = tx.QueryRowContext(ctx, `SELECT id FROM digital_keys WHERE offer_id=$1 AND order_item_id IS NULL ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED`, line.OfferID).Scan(&keyID)
				if err != nil {
					httpjson.Error(w, 409, "Stock changed during checkout; retry")
					return
				}
				_, err = tx.ExecContext(ctx, `UPDATE digital_keys SET order_item_id=$1,delivered_at=NOW() WHERE id=$2 AND order_item_id IS NULL`, itemID, keyID)
				if err != nil {
					httpjson.Error(w, 500, "Key delivery unavailable")
					return
				}
			}
		}
	}
	order, err := loadOrder(ctx, tx, id, u.ID)
	if err != nil {
		log.Printf("load checkout: %v", err)
		httpjson.Error(w, 500, "Could not load order")
		return
	}
	if err = tx.Commit(); err != nil {
		log.Printf("commit checkout: %v", err)
		httpjson.Error(w, 500, "Checkout unavailable")
		return
	}
	httpjson.Write(w, 201, order)
}
func loadOrder(ctx context.Context, db interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, id, userID int64) (Order, error) {
	var o Order
	o.Items = make([]PurchasedItem, 0)
	err := db.QueryRowContext(ctx, `SELECT id,status,total_cents,created_at FROM orders WHERE id=$1 AND user_id=$2`, id, userID).Scan(&o.ID, &o.Status, &o.TotalCents, &o.CreatedAt)
	if err != nil {
		return o, err
	}
	rows, err := db.QueryContext(ctx, `SELECT oi.id,oi.offer_id,p.title,p.slug,s.name,oi.price_cents,k.key_code FROM order_items oi JOIN offers off ON off.id=oi.offer_id JOIN products p ON p.id=off.product_id JOIN sellers s ON s.id=off.seller_id LEFT JOIN digital_keys k ON k.order_item_id=oi.id WHERE oi.order_id=$1 ORDER BY oi.id`, id)
	if err != nil {
		return o, err
	}
	defer rows.Close()
	for rows.Next() {
		var item PurchasedItem
		if err := rows.Scan(&item.ID, &item.OfferID, &item.Title, &item.Slug, &item.Seller, &item.PriceCents, &item.KeyCode); err != nil {
			return o, err
		}
		o.Items = append(o.Items, item)
	}
	return o, rows.Err()
}
func (s Service) History(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFrom(r.Context())
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	rows, err := s.DB.QueryContext(ctx, `SELECT id FROM orders WHERE user_id=$1 AND status <> 'pending' ORDER BY id DESC LIMIT 50`, u.ID)
	if err != nil {
		httpjson.Error(w, 500, "Orders unavailable")
		return
	}
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if e := rows.Scan(&id); e != nil {
			rows.Close()
			httpjson.Error(w, 500, "Orders unavailable")
			return
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		httpjson.Error(w, 500, "Orders unavailable")
		return
	}
	result := make([]Order, 0, len(ids))
	for _, id := range ids {
		o, e := loadOrder(ctx, s.DB, id, u.ID)
		if e != nil {
			httpjson.Error(w, 500, "Orders unavailable")
			return
		}
		result = append(result, o)
	}
	httpjson.Write(w, 200, map[string]any{"items": result})
}
func (s Service) Detail(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFrom(r.Context())
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.Error(w, 400, "Invalid order ID")
		return
	}
	o, err := loadOrder(r.Context(), s.DB, id, u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 404, "Order not found")
		return
	}
	if err != nil {
		httpjson.Error(w, 500, "Order unavailable")
		return
	}
	httpjson.Write(w, 200, o)
}
