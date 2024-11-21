package bgweb

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"vezgammon/server/config"
	"vezgammon/server/types"
)

var url = ""

func Init(conf *config.Config) {
	url = conf.Bgweb.Url
}

type Board struct {
	O CheckerLayout `json:"o"`
	X CheckerLayout `json:"x"`
}

type EngineConfig struct {
	MaxMoves   int
	ScoreMoves bool
}

var all_legal_moves_config = EngineConfig{
	MaxMoves:   100,
	ScoreMoves: false,
}

var get_best_move_config = EngineConfig{
	MaxMoves:   1,
	ScoreMoves: true,
}

type CheckerLayout struct {
	N1  int8 `json:"1"`
	N2  int8 `json:"2"`
	N3  int8 `json:"3"`
	N4  int8 `json:"4"`
	N5  int8 `json:"5"`
	N6  int8 `json:"6"`
	N7  int8 `json:"7"`
	N8  int8 `json:"8"`
	N9  int8 `json:"9"`
	N10 int8 `json:"10"`
	N11 int8 `json:"11"`
	N12 int8 `json:"12"`
	N13 int8 `json:"13"`
	N14 int8 `json:"14"`
	N15 int8 `json:"15"`
	N16 int8 `json:"16"`
	N17 int8 `json:"17"`
	N18 int8 `json:"18"`
	N19 int8 `json:"19"`
	N20 int8 `json:"20"`
	N21 int8 `json:"21"`
	N22 int8 `json:"22"`
	N23 int8 `json:"23"`
	N24 int8 `json:"24"`
	Bar int8 `json:"bar"`
}

type CheckerPlay struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type EvalInfo struct {
	Cubeful bool `json:"cubeful"`
	Plies   int  `json:"plies"`
}

type Evaluation struct {
	Diff        float32     `json:"diff"`
	Eq          float32     `json:"eq"`
	Info        EvalInfo    `json:"info"`
	Probability Probability `json:"probability"`
}

type Move struct {
	Evaluation Evaluation    `json:"evaluation"`
	Play       []CheckerPlay `json:"play"`
}

type MoveArgs struct {
	Board      Board  `json:"board"`
	Cubeful    bool   `json:"cubeful"`
	Dice       [2]int `json:"dice"`
	MaxMoves   int    `json:"max-moves"`
	Player     string `json:"player"`
	ScoreMoves bool   `json:"score-moves"`
}

type Probability struct {
	Lose   float32 `json:"lose"`
	LoseBG float32 `json:"loseBG"`
	LoseG  float32 `json:"loseG"`
	Win    float32 `json:"win"`
	WinBG  float32 `json:"winBG"`
	WinG   float32 `json:"winG"`
}

func (b *Board) toGame() *types.Game {
	var g types.Game

	g.P1Checkers[0] = b.O.Bar
	g.P1Checkers[1] = b.O.N1
	g.P1Checkers[2] = b.O.N2
	g.P1Checkers[3] = b.O.N3
	g.P1Checkers[4] = b.O.N4
	g.P1Checkers[5] = b.O.N5
	g.P1Checkers[6] = b.O.N6
	g.P1Checkers[7] = b.O.N7
	g.P1Checkers[8] = b.O.N8
	g.P1Checkers[9] = b.O.N9
	g.P1Checkers[10] = b.O.N10
	g.P1Checkers[11] = b.O.N11
	g.P1Checkers[12] = b.O.N12
	g.P1Checkers[13] = b.O.N13
	g.P1Checkers[14] = b.O.N14
	g.P1Checkers[15] = b.O.N15
	g.P1Checkers[16] = b.O.N16
	g.P1Checkers[17] = b.O.N17
	g.P1Checkers[18] = b.O.N18
	g.P1Checkers[19] = b.O.N19
	g.P1Checkers[20] = b.O.N20
	g.P1Checkers[21] = b.O.N21
	g.P1Checkers[22] = b.O.N22
	g.P1Checkers[23] = b.O.N23
	g.P1Checkers[24] = b.O.N24

	g.P2Checkers[0] = b.X.Bar
	g.P2Checkers[1] = b.X.N1
	g.P2Checkers[2] = b.X.N2
	g.P2Checkers[3] = b.X.N3
	g.P2Checkers[4] = b.X.N4
	g.P2Checkers[5] = b.X.N5
	g.P2Checkers[6] = b.X.N6
	g.P2Checkers[7] = b.X.N7
	g.P2Checkers[8] = b.X.N8
	g.P2Checkers[9] = b.X.N9
	g.P2Checkers[10] = b.X.N10
	g.P2Checkers[11] = b.X.N11
	g.P2Checkers[12] = b.X.N12
	g.P2Checkers[13] = b.X.N13
	g.P2Checkers[14] = b.X.N14
	g.P2Checkers[15] = b.X.N15
	g.P2Checkers[16] = b.X.N16
	g.P2Checkers[17] = b.X.N17
	g.P2Checkers[18] = b.X.N18
	g.P2Checkers[19] = b.X.N19
	g.P2Checkers[20] = b.X.N20
	g.P2Checkers[21] = b.X.N21
	g.P2Checkers[22] = b.X.N22
	g.P2Checkers[23] = b.X.N23
	g.P2Checkers[24] = b.X.N24

	return &g
}

func GametoMoveArgs(g *types.Game, engine_config EngineConfig) *MoveArgs {
	var moveargs MoveArgs

	var b Board

	b.O.Bar = g.P1Checkers[0]
	b.O.N1 = g.P1Checkers[1]
	b.O.N2 = g.P1Checkers[2]
	b.O.N3 = g.P1Checkers[3]
	b.O.N4 = g.P1Checkers[4]
	b.O.N5 = g.P1Checkers[5]
	b.O.N6 = g.P1Checkers[6]
	b.O.N7 = g.P1Checkers[7]
	b.O.N8 = g.P1Checkers[8]
	b.O.N9 = g.P1Checkers[9]
	b.O.N10 = g.P1Checkers[10]
	b.O.N11 = g.P1Checkers[11]
	b.O.N12 = g.P1Checkers[12]
	b.O.N13 = g.P1Checkers[13]
	b.O.N14 = g.P1Checkers[14]
	b.O.N15 = g.P1Checkers[15]
	b.O.N16 = g.P1Checkers[16]
	b.O.N17 = g.P1Checkers[17]
	b.O.N18 = g.P1Checkers[18]
	b.O.N19 = g.P1Checkers[19]
	b.O.N20 = g.P1Checkers[20]
	b.O.N21 = g.P1Checkers[21]
	b.O.N22 = g.P1Checkers[22]
	b.O.N23 = g.P1Checkers[23]
	b.O.N24 = g.P1Checkers[24]

	b.X.Bar = g.P2Checkers[0]
	b.X.N1 = g.P2Checkers[1]
	b.X.N2 = g.P2Checkers[2]
	b.X.N3 = g.P2Checkers[3]
	b.X.N4 = g.P2Checkers[4]
	b.X.N5 = g.P2Checkers[5]
	b.X.N6 = g.P2Checkers[6]
	b.X.N7 = g.P2Checkers[7]
	b.X.N8 = g.P2Checkers[8]
	b.X.N9 = g.P2Checkers[9]
	b.X.N10 = g.P2Checkers[10]
	b.X.N11 = g.P2Checkers[11]
	b.X.N12 = g.P2Checkers[12]
	b.X.N13 = g.P2Checkers[13]
	b.X.N14 = g.P2Checkers[14]
	b.X.N15 = g.P2Checkers[15]
	b.X.N16 = g.P2Checkers[16]
	b.X.N17 = g.P2Checkers[17]
	b.X.N18 = g.P2Checkers[18]
	b.X.N19 = g.P2Checkers[19]
	b.X.N20 = g.P2Checkers[20]
	b.X.N21 = g.P2Checkers[21]
	b.X.N22 = g.P2Checkers[22]
	b.X.N23 = g.P2Checkers[23]
	b.X.N24 = g.P2Checkers[24]

	var p string
	if g.CurrentPlayer == types.GameCurrentPlayerP1 {
		p = "o"
	} else {
		p = "x"
	}

	moveargs.Board = b
	moveargs.Cubeful = true // always play with cube
	moveargs.Dice = g.Dices
	moveargs.Player = p
	moveargs.MaxMoves = engine_config.MaxMoves
	moveargs.ScoreMoves = engine_config.ScoreMoves

	return &moveargs
}

// dont't set Dices, had to be done separately
func (m *Move) toTurn() (*types.Turn, error) {
	var t types.Turn
	var err error
	for _, play := range m.Play {
		var to, from int64
		if play.From == "bar" {
			from = 0
		} else {
			from, err = strconv.ParseInt(play.From, 10, 64)
		}

		if play.To == "off" {
			to = 25
		} else {
			to, err = strconv.ParseInt(play.To, 10, 64)
		}

		m := types.Move{From: from, To: to}
		t.Moves = append(t.Moves, m)
	}

	t.Double = false // engine can't double
	return &t, err
}

var DefaultMoveArgs MoveArgs = MoveArgs{
	Board:      Board{},
	Cubeful:    false,
	ScoreMoves: false, // get list of legal moves
	MaxMoves:   100,   // get all moves
}

func GetMoves(moveargs *MoveArgs) ([]Move, error) {

	apiurl := url + "getmoves"

	postbody, err := json.Marshal(moveargs)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://"+apiurl, "application/json", bytes.NewReader(postbody))
	if err != nil {
		return nil, err
	}

	var m []Move
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buff, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func MoveArrayToMoveArrayArray(movesarray []Move) [][]types.Move {
	retarrayarray := make([][]types.Move, 0)

	for _, checkerplayarray := range movesarray {
		var retarray []types.Move
		for _, play := range checkerplayarray.Play {
			var to, from int64
			if play.From == "bar" {
				from = 0
			} else {
				from, _ = strconv.ParseInt(play.From, 10, 64)
			}

			if play.To == "off" {
				to = 25
			} else {
				to, _ = strconv.ParseInt(play.To, 10, 64)
			}

			retarray = append(retarray, types.Move{From: from, To: to})
		}
		retarrayarray = append(retarrayarray, retarray)
	}
	return retarrayarray
}

func GetLegalMoves(g *types.Game) ([][]types.Move, error) {
	mv := GametoMoveArgs(g, all_legal_moves_config)

	// slog.With("moves args", *mv).Debug("Game to move args")

	moves1, err := GetMoves(mv)
	if err != nil {
		return nil, err
	}

	mv.Dice[0], mv.Dice[1] = mv.Dice[1], mv.Dice[0]

	moves2, err := GetMoves(mv)
	if err != nil {
		return nil, err
	}

	moves := append(moves1, moves2...)

	// slog.With("moves", moves).Debug("Got moves")

	possibleMoves := MoveArrayToMoveArrayArray(moves)

	// Dio perdonami
	modified := true
	for modified {
		possibleMoves, modified = fill(possibleMoves, g)
	}

	if len(possibleMoves) == 0 {
		possibleMoves = make([][]types.Move, 0)
	}

	return possibleMoves, nil
}

func GetBestMove(g *types.Game) (*types.Turn, error) {
	mv := GametoMoveArgs(g, get_best_move_config)

	moves, err := GetMoves(mv)
	if err != nil {
		return nil, err
	}

	turn, err := moves[0].toTurn()
	if err != nil {
		return nil, err
	}

	return turn, nil
}

func GetEasyMove(g *types.Game) (*types.Turn, error) {
	conf := EngineConfig{
		MaxMoves:   5,
		ScoreMoves: true,
	}
	mv := GametoMoveArgs(g, conf)

	moves, err := GetMoves(mv)
	if err != nil {
		return nil, err
	}

	move := moves[len(moves)-1] // get second best move

	turn, err := move.toTurn()
	if err != nil {
		return nil, err
	}

	return turn, nil
}

func GetMediumMove(g *types.Game) (*types.Turn, error) {
	conf := EngineConfig{
		MaxMoves:   3,
		ScoreMoves: true,
	}
	mv := GametoMoveArgs(g, conf)

	moves, err := GetMoves(mv)
	if err != nil {
		return nil, err
	}

	move := moves[len(moves)-1] // get third best move

	turn, err := move.toTurn()
	if err != nil {
		return nil, err
	}

	return turn, nil
}
