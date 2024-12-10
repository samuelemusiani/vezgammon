package bgweb

import (
	"encoding/json"
	"testing"
	"vezgammon/server/types"

	"gotest.tools/v3/assert"
)

func TestGetmoves(t *testing.T) {

	url = "localhost:3030/api/v1/"

	jsonArg := `
 {
   "board": {
     "o": {
       "6": 5,
       "8": 3,
       "13": 5,
       "24": 2
     },
     "x": {
       "6": 5,
       "8": 3,
       "13": 5,
       "24": 2
     }
   },
   "cubeful": false,
   "dice": [
     3,
     1
   ],
   "max-moves": 3,
   "player": "x",
   "score-moves": true
 }
 	`

	jsonresp := `
   [
   {
     "evaluation": {
       "diff": 0,
       "eq": 0.159,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.449,
         "loseBG": 0.005,
         "loseG": 0.124,
         "win": 0.551,
         "winBG": 0.013,
         "winG": 0.174
       }
     },
     "play": [
       {
         "from": "8",
         "to": "5"
       },
       {
         "from": "6",
         "to": "5"
       }
     ]
   },
   {
     "evaluation": {
       "diff": -0.168,
       "eq": -0.009,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.503,
         "loseBG": 0.007,
         "loseG": 0.14,
         "win": 0.497,
         "winBG": 0.008,
         "winG": 0.137
       }
     },
     "play": [
       {
         "from": "13",
         "to": "10"
       },
       {
         "from": "24",
         "to": "23"
       }
     ]
   },
   {
     "evaluation": {
       "diff": -0.175,
       "eq": -0.015,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.503,
         "loseBG": 0.004,
         "loseG": 0.135,
         "win": 0.497,
         "winBG": 0.005,
         "winG": 0.125
       }
     },
     "play": [
       {
         "from": "24",
         "to": "21"
       },
       {
         "from": "21",
         "to": "20"
       }
     ]
   }
 ]
   `

	moveargs := DefaultMoveArgs
	err := json.Unmarshal([]byte(jsonArg), &moveargs)
	assert.NilError(t, err)
	t.Logf("moveargs %v", moveargs)

	var respmoves []Move
	err = json.Unmarshal([]byte(jsonresp), &respmoves)
	assert.NilError(t, err)
	t.Logf("respmoves %v", respmoves)

	rp, err := GetMoves(&moveargs)
	assert.NilError(t, err)

	assert.DeepEqual(t, rp, respmoves)
}

func TestBoardToGame(t *testing.T) {
	board := Board{
		O: CheckerLayout{
			N1:  5,
			N2:  3,
			Bar: 0,
		},
		X: CheckerLayout{
			N1:  2,
			N2:  4,
			Bar: 1,
		},
	}

	game := board.toGame()
	assert.Equal(t, board.O.N1, game.P1Checkers[1])
	assert.Equal(t, board.O.N2, game.P1Checkers[2])
	assert.Equal(t, board.O.Bar, game.P1Checkers[0])
	assert.Equal(t, board.X.N1, game.P2Checkers[1])
	assert.Equal(t, board.X.N2, game.P2Checkers[2])
	assert.Equal(t, board.X.Bar, game.P2Checkers[0])
}

func TestGametoMoveArgs(t *testing.T) {
	var game types.Game
	game.P1Checkers[0] = 0
	game.P1Checkers[1] = 5
	game.P2Checkers[1] = 3

	game.CurrentPlayer = types.GameCurrentPlayerP1
	game.Dices = types.Dices{3, 1}

	engineconfig := EngineConfig{
		MaxMoves:   3,
		ScoreMoves: true,
	}

	moveargs := GametoMoveArgs(&game, engineconfig)

	assert.Equal(t, moveargs.Player, "o")
	assert.DeepEqual(t, moveargs.Dice, [2]int{3, 1})
	assert.Equal(t, moveargs.MaxMoves, engineconfig.MaxMoves)
	assert.Equal(t, moveargs.ScoreMoves, engineconfig.ScoreMoves)
	assert.Equal(t, moveargs.Cubeful, true)

	assert.Equal(t, moveargs.Board.O.N1, game.P1Checkers[1])
	assert.Equal(t, moveargs.Board.X.N1, game.P2Checkers[1])
}

func TestMoveArrayToMoveArrayArray(t *testing.T) {
	movesarray := []Move{
		{
			Play: []CheckerPlay{
				{
					From: "8",
					To:   "5",
				},
				{
					From: "6",
					To:   "5",
				},
			},
		},
		{
			Play: []CheckerPlay{
				{
					From: "13",
					To:   "10",
				},
				{
					From: "24",
					To:   "23",
				},
			},
		},
	}

	typemovesarray := [][]types.Move{
		{
			{
				From: 8,
				To:   5,
			},
			{
				From: 6,
				To:   5,
			},
		},
		{
			{
				From: 13,
				To:   10,
			},
			{
				From: 24,
				To:   23,
			},
		},
	}

	assert.DeepEqual(t, MoveArrayToMoveArrayArray(movesarray), typemovesarray)
}

func TestGetLegalMoves(t *testing.T) {
	typemovesarray := [][]types.Move{
		{{From: 24, To: 21}, {From: 24, To: 23}},
		{{From: 24, To: 21}, {From: 21, To: 20}},
		{{From: 24, To: 21}, {From: 8, To: 7}},
		{{From: 24, To: 21}, {From: 6, To: 5}},
		{{From: 13, To: 10}, {From: 24, To: 23}},
		{{From: 13, To: 10}, {From: 10, To: 9}},
		{{From: 13, To: 10}, {From: 8, To: 7}},
		{{From: 13, To: 10}, {From: 6, To: 5}},
		{{From: 8, To: 5}, {From: 24, To: 23}},
		{{From: 8, To: 5}, {From: 8, To: 7}},
		{{From: 8, To: 5}, {From: 6, To: 5}},
		{{From: 8, To: 5}, {From: 5, To: 4}},
		{{From: 6, To: 3}, {From: 24, To: 23}},
		{{From: 6, To: 3}, {From: 8, To: 7}},
		{{From: 6, To: 3}, {From: 6, To: 5}},
		{{From: 6, To: 3}, {From: 3, To: 2}},
		{{From: 24, To: 23}, {From: 24, To: 21}},
		{{From: 24, To: 23}, {From: 23, To: 20}},
		{{From: 24, To: 23}, {From: 13, To: 10}},
		{{From: 24, To: 23}, {From: 8, To: 5}},
		{{From: 24, To: 23}, {From: 6, To: 3}},
		{{From: 8, To: 7}, {From: 24, To: 21}},
		{{From: 8, To: 7}, {From: 13, To: 10}},
		{{From: 8, To: 7}, {From: 8, To: 5}},
		{{From: 8, To: 7}, {From: 7, To: 4}},
		{{From: 8, To: 7}, {From: 6, To: 3}},
		{{From: 6, To: 5}, {From: 24, To: 21}},
		{{From: 6, To: 5}, {From: 13, To: 10}},
		{{From: 6, To: 5}, {From: 8, To: 5}},
		{{From: 6, To: 5}, {From: 6, To: 3}},
		{{From: 6, To: 5}, {From: 5, To: 2}},
		{{From: 13, To: 10}, {From: 10, To: 9}},
	}

	game := types.Game{
		P1Checkers:    [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		P2Checkers:    [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		CurrentPlayer: types.GameCurrentPlayerP1,
		Dices:         types.Dices{3, 1},
	}

	legalmoves, err := GetLegalMoves(&game)
	assert.NilError(t, err)
	assert.DeepEqual(t, legalmoves, typemovesarray)
}
