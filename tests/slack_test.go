package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goshanraj-g/vigil/internal/slack"
)

func TestSlashCommand_EmptyText(t *testing.T) {
	h := slack.New(nil)

	req := httptest.NewRequest(http.MethodPost, "/slack/command", strings.NewReader("text="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.SlashCommand(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Usage") {
		t.Errorf("expected usage message, got: %s", w.Body.String())
	}
}

func TestSlashCommand_StopMissingNumber(t *testing.T) {
	h := slack.New(nil)

	req := httptest.NewRequest(http.MethodPost, "/slack/command", strings.NewReader("text=stop"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.SlashCommand(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Usage") {
		t.Errorf("expected usage message, got: %s", w.Body.String())
	}
}

func TestSlashCommand_StopInvalidNumber(t *testing.T) {
	h := slack.New(nil)

	req := httptest.NewRequest(http.MethodPost, "/slack/command", strings.NewReader("text=stop+abc"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.SlashCommand(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invalid number") {
		t.Errorf("expected invalid number message, got: %s", w.Body.String())
	}
}
