package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/goshanraj-g/vigil/internal/search"
)

// mockSerper starts a fake Serper server and returns a search client pointed at it
func mockSerper(t *testing.T, response any) (*search.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	os.Setenv("SERPER_API_KEY", "test-key")
	c, err := search.NewWithURL(srv.URL + "/search")
	if err != nil {
		t.Fatalf("search.NewWithURL: %v", err)
	}
	return c, srv
}

func TestSearch_ReturnsResults(t *testing.T) {
	fakeResponse := map[string]any{
		"organic": []map[string]any{
			{"link": "https://example.com/article", "title": "Test Article", "snippet": "This is a test."},
		},
	}
	client, srv := mockSerper(t, fakeResponse)
	defer srv.Close()

	results, err := client.Search(t.Context(), "test query")
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].URL != "https://example.com/article" {
		t.Errorf("unexpected URL: %s", results[0].URL)
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	os.Setenv("SERPER_API_KEY", "test-key")
	c, err := search.NewWithURL("http://localhost")
	if err != nil {
		t.Fatalf("search.NewWithURL: %v", err)
	}
	_, err = c.Search(t.Context(), "")
	if err == nil {
		t.Error("expected error for empty query, got nil")
	}
}

func TestSearch_SkipsEmptyLinks(t *testing.T) {
	fakeResponse := map[string]any{
		"organic": []map[string]any{
			{"link": "", "title": "No Link", "snippet": "Missing URL."},
			{"link": "https://example.com", "title": "Valid", "snippet": "Has URL."},
		},
	}
	client, srv := mockSerper(t, fakeResponse)
	defer srv.Close()

	results, err := client.Search(t.Context(), "test")
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (empty link skipped), got %d", len(results))
	}
}
