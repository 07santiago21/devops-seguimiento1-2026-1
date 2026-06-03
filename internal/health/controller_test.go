package health

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReturnsHealthPayload(t *testing.T) {
	handler := Handler(Config{
		Healthcheck: "canary",
		Version:     "1.2.0",
		DeployedAt:  "2026-06-02T10:00:00Z",
		Service:     "students-api",
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}

	body := w.Body.Bytes()
	for _, want := range [][]byte{
		[]byte(`"status":"ok"`),
		[]byte(`"healthcheck":"canary"`),
		[]byte(`"version":"1.2.0"`),
		[]byte(`"deployed_at":"2026-06-02T10:00:00Z"`),
		[]byte(`"service":"students-api"`),
	} {
		if !bytes.Contains(body, want) {
			t.Fatalf("expected response to contain %s, got %s", string(want), w.Body.String())
		}
	}
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	handler := Handler(Config{})

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
