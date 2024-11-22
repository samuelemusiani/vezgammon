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

func TestLogin(t *testing.T) {
	loginuser := loginUserType{
		Username: "testregisteruser",
		Password: "1234",
	}

	expectedresponse := loginResponseType{
		Message: "Login successful",
		User: loginResponseUser{
			Username: "testregisteruser",
			Email:    "a@a.it",
		},
	}

	w := httptest.NewRecorder()
	loginuserjson, _ := json.Marshal(loginuser)
	req, err := http.NewRequest("POST", "/api/login", strings.NewReader(string(loginuserjson)))
	assert.NilError(t, err)

	router.ServeHTTP(w, req)

	var retresponse loginResponseType
	err = json.Unmarshal(w.Body.Bytes(), &retresponse)
	assert.NilError(t, err)

	session_token = w.Result().Cookies()[0]

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, expectedresponse.Message, retresponse.Message)
	assert.Equal(t, expectedresponse.User.Username, retresponse.User.Username)
	assert.Equal(t, expectedresponse.User.Email, retresponse.User.Email)
}

func TestLogout(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/logout", nil)
	assert.NilError(t, err)
	req.AddCookie(session_token)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}
