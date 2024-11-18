package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vezgammon/server/db"

	"gotest.tools/v3/assert"
)

func TestStartGameLocalcally(t *testing.T) {
	q := "DROP TABLE IF EXISTS games CASCADE"
	_, err := db.Conn.Exec(q)
	assert.NilError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/play/local", nil)
	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
	// ...
}
