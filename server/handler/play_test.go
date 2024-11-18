package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartGameLocalcally(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/play/local", nil)
	router.ServeHTTP(w, req)
	// ...
}
