package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"eneba-marketplace/backend/internal/api"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMarketplaceJourney(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("Set TEST_DATABASE_URL for isolated PostgreSQL integration tests")
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	handler := api.NewRouter(db)
	req := func(method, path string, body any, cookie *http.Cookie) *httptest.ResponseRecorder {
		var data []byte
		if body != nil {
			data, _ = json.Marshal(body)
		}
		r := httptest.NewRequest(method, path, bytes.NewReader(data))
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	email := fmt.Sprintf("test-%d@example.test", time.Now().UnixNano())
	register := req("POST", "/api/auth/register", map[string]any{"email": email, "password": "unit-test-password123"}, nil)
	if register.Code != 201 {
		t.Fatalf("register: %d %s", register.Code, register.Body.String())
	}
	var cookie *http.Cookie
	for _, c := range register.Result().Cookies() {
		if c.Name == "marketplace_session" {
			cookie = c
		}
	}
	if cookie == nil || !cookie.HttpOnly {
		t.Fatal("missing HttpOnly session cookie")
	}
	if w := req("GET", "/api/orders", nil, nil); w.Code != 401 {
		t.Fatalf("unauthorized history status %d", w.Code)
	}
	if w := req("GET", "/api/products?q=minecrft", nil, nil); w.Code != 200 || !strings.Contains(w.Body.String(), "Minecraft") {
		t.Fatalf("typo search: %d %s", w.Code, w.Body.String())
	}
	if w := req("GET", "/api/products?platform=PC&region=Europe&sort=price_asc&page=1&limit=5", nil, nil); w.Code != 200 {
		t.Fatalf("filters: %s", w.Body.String())
	}
	var offerID int64
	err = db.QueryRow(`SELECT o.id FROM offers o WHERE (SELECT COUNT(*) FROM digital_keys k WHERE k.offer_id=o.id AND k.order_item_id IS NULL)>=2 ORDER BY o.id DESC LIMIT 1`).Scan(&offerID)
	if err != nil {
		t.Fatal(err)
	}
	purchase := map[string]any{"items": []any{map[string]any{"offer_id": offerID, "quantity": 1}}, "payment_method": "approve", "idempotency_key": "integration-paid-0001"}
	paid := req("POST", "/api/checkout", purchase, cookie)
	if paid.Code != 201 {
		t.Fatalf("purchase: %d %s", paid.Code, paid.Body.String())
	}
	var order struct {
		ID    int64 `json:"id"`
		Items []struct {
			KeyCode *string `json:"key_code"`
		} `json:"items"`
	}
	if err = json.Unmarshal(paid.Body.Bytes(), &order); err != nil {
		t.Fatal(err)
	}
	if order.ID <= 0 || len(order.Items) != 1 || order.Items[0].KeyCode == nil {
		t.Fatalf("no delivered key: %s", paid.Body.String())
	}
	same := req("POST", "/api/checkout", purchase, cookie)
	if same.Code != 200 {
		t.Fatalf("retry did not reuse order: %d %s", same.Code, same.Body.String())
	}
	var replay struct {
		ID int64 `json:"id"`
	}
	_ = json.Unmarshal(same.Body.Bytes(), &replay)
	if replay.ID != order.ID {
		t.Fatalf("duplicate order: %d vs %d", order.ID, replay.ID)
	}
	purchase["payment_method"] = "decline"
	conflict := req("POST", "/api/checkout", purchase, cookie)
	if conflict.Code != 409 {
		t.Fatalf("modified idempotent request: %d", conflict.Code)
	}
	purchase["idempotency_key"] = "integration-declined-0001"
	declined := req("POST", "/api/checkout", purchase, cookie)
	if declined.Code != 201 || !strings.Contains(declined.Body.String(), `"status":"failed"`) || strings.Contains(declined.Body.String(), "DEMO-ONLY-") {
		t.Fatalf("declined behavior: %d %s", declined.Code, declined.Body.String())
	}
	history := req("GET", "/api/orders", nil, cookie)
	if history.Code != 200 || !strings.Contains(history.Body.String(), "DEMO-ONLY-") {
		t.Fatalf("history: %d %s", history.Code, history.Body.String())
	}
	if w := req("GET", fmt.Sprintf("/api/orders/%d", order.ID), nil, nil); w.Code != 401 {
		t.Fatalf("private order unprotected: %d", w.Code)
	}
	second := req("POST", "/api/auth/register", map[string]any{"email": fmt.Sprintf("other-%d@example.test", time.Now().UnixNano()), "password": "unit-test-password123"}, nil)
	if second.Code != 201 {
		t.Fatalf("second user: %d %s", second.Code, second.Body.String())
	}
	otherCookie := second.Result().Cookies()[0]
	if w := req("GET", fmt.Sprintf("/api/orders/%d", order.ID), nil, otherCookie); w.Code != 404 {
		t.Fatalf("cross-user order visible: %d", w.Code)
	}
	if w := req("POST", "/api/auth/logout", map[string]any{}, cookie); w.Code != 200 {
		t.Fatalf("logout: %d", w.Code)
	}
	if w := req("GET", "/api/orders", nil, cookie); w.Code != 401 {
		t.Fatalf("logout did not revoke session: %d", w.Code)
	}
}
