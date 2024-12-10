package ws

import (
	"errors"
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

type MockDB struct{}

var errCurrentGame = false
var errGetGame = false

func (m *MockDB) GetCurrentGame(userID int64) (*types.ReturnGame, error) {
	if errCurrentGame {
		return nil, errors.New("error")
	}
	return &types.ReturnGame{
		Player1: "1",
		Player2: "2",
	}, nil
}

func (m *MockDB) GetGame(gameID int64) (*types.Game, error) {
	if errGetGame {
		return nil, errors.New("error")
	}
	return &types.Game{
		Player1: 1,
		Player2: 2,
	}, nil
}

func TestChatRespondeToMessage(t *testing.T) {
	Init(&MockDB{})
	err := chatRespondeToMessage(1, Message{})
	assert.ErrorIs(t, err, ErrWritingMessage)

	err = chatRespondeToMessage(22, Message{})
	assert.ErrorIs(t, err, ErrWritingMessage)
}

func TestChatRespondeToMessageError(t *testing.T) {
	Init(&MockDB{})
	errCurrentGame = true
	chatRespondeToMessage(1, Message{})

	errCurrentGame = false
	errGetGame = true
	chatRespondeToMessage(1, Message{})
}
