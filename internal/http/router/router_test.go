package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/router"
)

func TestHealthz(t *testing.T) {
	handler := router.New(nil, nil)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	expectedBody := `{"status":"ok"}`
	if response.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, response.Body.String())
	}
}
