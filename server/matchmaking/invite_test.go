package matchmaking

import (
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

type MockDB struct{}

func (m *MockDB) GetUser(id int64) (*types.User, error) {
	switch id {
	case 1:
		return &types.User{ID: 1, Username: "User1"}, nil
	case 2:
		return &types.User{ID: 2, Username: "User2"}, nil
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

func TestGenerateLink(t *testing.T) {
	link, err := GenerateLink(1)
	assert.NilError(t, err)
	assert.Assert(t, link != "")

	suuid, ok := users[1]
	assert.Assert(t, ok)
	assert.Equal(t, suuid, link)

	id, ok := links[link]
	assert.Assert(t, ok)
	assert.Equal(t, id, int64(1))
}

func TestJoinLink(t *testing.T) {
	db = &MockDB{}
	ws = &MockWS{}

	link, err := GenerateLink(1)
	assert.NilError(t, err)

	err = JoinLink(link, 2)
	assert.NilError(t, err)

	_, ok := users[1]
	assert.Assert(t, !ok)

	_, ok = links[link]
	assert.Assert(t, !ok)
}
