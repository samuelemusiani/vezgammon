package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder()

	testRegisterUser := customUser{
		Password:  "1234",
		Username:  "testregisteruser",
		Firstname: "testregisteruser",
		Lastname:  "testregisteruser",
		Mail:      "a@a.it",
	}

	user := types.User{
		Username:  "testregisteruser",
		Firstname: "testregisteruser",
		Lastname:  "testregisteruser",
		Mail:      "a@a.it",
	}

	userjson, _ := json.Marshal(testRegisterUser)

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(string(userjson)))
	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	var retuser types.User
	err = json.Unmarshal(w.Body.Bytes(), &retuser)
	assert.NilError(t, err)

	assert.Equal(t, user.Username, retuser.Username)
	assert.Equal(t, user.Firstname, retuser.Firstname)
	assert.Equal(t, user.Lastname, retuser.Lastname)
	assert.Equal(t, user.Mail, retuser.Mail)
	assert.Equal(t, w.Code, http.StatusCreated)
}
