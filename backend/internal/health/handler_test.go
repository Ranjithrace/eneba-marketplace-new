package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePinger struct{ err error }

func (f fakePinger) PingContext(context.Context) error { return f.err }

func TestHealthy(t *testing.T) {
	w := httptest.NewRecorder()
	Handler(fakePinger{}).ServeHTTP(w, httptest.NewRequest("GET", "/api/health", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"database":"connected"`) {
		t.Fatalf("unexpected response %d %s", w.Code, w.Body.String())
	}
}
func TestUnavailable(t *testing.T) {
	w := httptest.NewRecorder()
	Handler(fakePinger{errors.New("db down")}).ServeHTTP(w, httptest.NewRequest("GET", "/api/health", nil))
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(w.Body.String(), `"database":"unavailable"`) {
		t.Fatalf("unexpected response %d %s", w.Code, w.Body.String())
	}
}
