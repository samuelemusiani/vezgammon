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

func TestGameSetup(t *testing.T) {
	tuserr := customUser{
		Password:  "1234",
		Username:  "testgameuser",
		Firstname: "testgameuser",
		Lastname:  "testgameuser",
		Mail:      "testgameuser@a.it",
	}

	tuserl := loginUserType{
		Username: "testgameuser",
		Password: "1234",
	}

	w := httptest.NewRecorder()

	// register
	userjson, _ := json.Marshal(tuserr)

	req, err := http.NewRequest("POST", "/api/register", strings.NewReader(string(userjson)))
	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)

	//login
	w = httptest.NewRecorder()
	loginuserjson, _ := json.Marshal(tuserl)
	req, err = http.NewRequest("POST", "/api/login", strings.NewReader(string(loginuserjson)))
	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	// get session token
	session_token = w.Result().Cookies()[0]
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestStartGameLocalcally(t *testing.T) {

	expectedresponse := types.NewGame{
		Game: types.ReturnGame{
			Player1:      "testgameuser",
			Player2:      "testgameuser",
			Status:       types.GameStatusOpen,
			P1Checkers:   [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			P2Checkers:   [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			DoubleValue:  1,
			DoubleOwner:  types.GameDoubleOwnerAll,
			WantToDouble: false,
		},
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/play/local", nil)
	assert.NilError(t, err)

	assert.Assert(t, session_token != nil)
	req.AddCookie(session_token)

	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	var retresponse types.NewGame
	err = json.Unmarshal(w.Body.Bytes(), &retresponse)
	assert.NilError(t, err)

	assert.Equal(t, w.Code, http.StatusCreated)

	assert.Equal(t, expectedresponse.Game.Player1, retresponse.Game.Player1)
	assert.Equal(t, expectedresponse.Game.Player2, retresponse.Game.Player2)
	assert.Equal(t, expectedresponse.Game.Status, retresponse.Game.Status)
	assert.DeepEqual(t, expectedresponse.Game.P1Checkers, retresponse.Game.P1Checkers)
	assert.DeepEqual(t, expectedresponse.Game.P2Checkers, retresponse.Game.P2Checkers)
	assert.Equal(t, expectedresponse.Game.DoubleValue, retresponse.Game.DoubleValue)
	assert.Equal(t, expectedresponse.Game.DoubleOwner, retresponse.Game.DoubleOwner)
	assert.Equal(t, expectedresponse.Game.WantToDouble, retresponse.Game.WantToDouble)
}

func TestGetCurrentGame(t *testing.T) {
	expectedresponse := types.ReturnGame{
		Player1:      "testgameuser",
		Player2:      "testgameuser",
		Status:       types.GameStatusOpen,
		P1Checkers:   [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		P2Checkers:   [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		DoubleValue:  1,
		DoubleOwner:  types.GameDoubleOwnerAll,
		WantToDouble: false,
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/play/", nil)
	assert.NilError(t, err)

	assert.Assert(t, session_token != nil)
	req.AddCookie(session_token)

	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	var retresponse types.ReturnGame
	err = json.Unmarshal(w.Body.Bytes(), &retresponse)
	assert.NilError(t, err)

	assert.Equal(t, w.Code, http.StatusOK)

	assert.Equal(t, expectedresponse.Player1, retresponse.Player1)
	assert.Equal(t, expectedresponse.Player2, retresponse.Player2)
	assert.Equal(t, expectedresponse.Status, retresponse.Status)
	assert.DeepEqual(t, expectedresponse.P1Checkers, retresponse.P1Checkers)
	assert.DeepEqual(t, expectedresponse.P2Checkers, retresponse.P2Checkers)
	assert.Equal(t, expectedresponse.DoubleValue, retresponse.DoubleValue)
	assert.Equal(t, expectedresponse.DoubleOwner, retresponse.DoubleOwner)
	assert.Equal(t, expectedresponse.WantToDouble, retresponse.WantToDouble)
}

// keep this test last
func TestSurrendToCurrentGame(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/api/play/", nil)
	assert.NilError(t, err)

	assert.Assert(t, session_token != nil)
	req.AddCookie(session_token)

	assert.NilError(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
}
