package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goshanraj-g/vigil/internal/api"
)

func TestCreateMonitor_BadJSON(t *testing.T) {
	r := api.NewRouter(nil)

	req := httptest.NewRequest(http.MethodPost, "/monitors", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateMonitor_MissingFields(t *testing.T) {
	r := api.NewRouter(nil)

	// missing webhook_url and schedule
	req := httptest.NewRequest(http.MethodPost, "/monitors", strings.NewReader(`{"query":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateMonitor_MissingQuery(t *testing.T) {
	r := api.NewRouter(nil)

	req := httptest.NewRequest(http.MethodPost, "/monitors", strings.NewReader(`{"webhook_url":"https://example.com","schedule":"@every 1m"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
