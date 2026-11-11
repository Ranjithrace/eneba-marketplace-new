package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRepo struct {
	response ProductsResponse
	err      error
}

func (f fakeRepo) List(_ context.Context, _ Filters) (ProductsResponse, error) {
	return f.response, f.err
}
func TestListHandlerProducts(t *testing.T) {
	p := int64(714)
	r := httptest.NewRecorder()
	ListHandler(fakeRepo{response: ProductsResponse{Items: []Product{{ID: 1, Title: "Minecraft: Java Edition", LowestPriceCents: &p}}, Count: 1, Total: 1, Page: 1, PageSize: 12, Pages: 1}}).ServeHTTP(r, httptest.NewRequest("GET", "/api/products", nil))
	if r.Code != 200 {
		t.Fatalf("status %d", r.Code)
	}
	var got ProductsResponse
	if e := json.Unmarshal(r.Body.Bytes(), &got); e != nil {
		t.Fatal(e)
	}
	if got.Total != 1 || got.Items[0].Title != "Minecraft: Java Edition" || *got.Items[0].LowestPriceCents != 714 {
		t.Fatalf("bad data: %+v", got)
	}
}
func TestInvalidFilters(t *testing.T) {
	for _, uri := range []string{"/api/products?page=0", "/api/products?limit=100", "/api/products?min_price=-1", "/api/products?max_price=hello"} {
		r := httptest.NewRecorder()
		ListHandler(fakeRepo{}).ServeHTTP(r, httptest.NewRequest("GET", uri, nil))
		if r.Code != 400 {
			t.Errorf("%s: expected 400, got %d", uri, r.Code)
		}
	}
}
func TestListHandlerMasksErrors(t *testing.T) {
	r := httptest.NewRecorder()
	ListHandler(fakeRepo{err: errors.New("DATABASE_PASSWORD=secret")}).ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/api/products", nil))
	if r.Code != 500 || strings.Contains(r.Body.String(), "secret") {
		t.Fatalf("error leak %d %s", r.Code, r.Body.String())
	}
}
