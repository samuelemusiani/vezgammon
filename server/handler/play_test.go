package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestStartGameLocalcally(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/play/local", nil)
	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	// assert.Equal(t, w.Code, http.StatusCreated)
	// ...
}
