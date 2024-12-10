package matchmaking

import (
	"database/sql"
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

type MockDB struct{}

func (m *MockDB) GetUser(id int64) (*types.User, error) {
	switch id {
	case 1:
		return &types.User{ID: 1, Username: "User1", Elo: 1000}, nil
	case 2:
		return &types.User{ID: 2, Username: "User2", Elo: 1000}, nil
	default:
		panic("Invalid user id on DB MOCK")
	}
}

func (m *MockDB) CreateGame(g types.Game) (*types.Game, error) {
	return nil, nil
}

type MockWS struct{}

func (m *MockWS) SendGameFound(int64) error {
	return nil
}

func initMock() {
	db = &MockDB{}
	ws = &MockWS{}
}

func TestCreateGame(t *testing.T) {
	initMock()
	_, err := CreateGame(1, 2, sql.NullInt64{Valid: false})
	assert.NilError(t, err)
}

func TestCheckIfValidOpponent(t *testing.T) {
	assert.Assert(t, checkIfValidOpponent(1000, 1000))
	assert.Assert(t, checkIfValidOpponent(1000, 1001))
	assert.Assert(t, !checkIfValidOpponent(1000, 1500))
	assert.Assert(t, !checkIfValidOpponent(1000, 500))
}

func TestMatchmaking(t *testing.T) {
	err := SearchGame(1)
	assert.NilError(t, err)
	err = SearchGame(2)
	assert.NilError(t, err)

	worker(true)
}

func TestStopSearch(t *testing.T) {
	err := SearchGame(1)
	assert.NilError(t, err)
	err = SearchGame(2)
	assert.NilError(t, err)

	StopSearch(1)
	StopSearch(2)

	assert.Assert(t, length() == 0)
}
