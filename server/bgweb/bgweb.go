package bgweb

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Board struct {
	O CheckerLayout `json:"o"`
	X CheckerLayout `json:"x"`
}

type CheckerLayout struct {
	N1  int `json:"1"`
	N2  int `json:"2"`
	N3  int `json:"3"`
	N4  int `json:"4"`
	N5  int `json:"5"`
	N6  int `json:"6"`
	N7  int `json:"7"`
	N8  int `json:"8"`
	N9  int `json:"9"`
	N10 int `json:"10"`
	N11 int `json:"11"`
	N12 int `json:"12"`
	N13 int `json:"13"`
	N14 int `json:"14"`
	N15 int `json:"15"`
	N16 int `json:"16"`
	N17 int `json:"17"`
	N18 int `json:"18"`
	N19 int `json:"19"`
	N20 int `json:"20"`
	N21 int `json:"21"`
	N22 int `json:"22"`
	N23 int `json:"23"`
	N24 int `json:"24"`
	Bar int `json:"bar"`
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

var DefaultMoveArgs MoveArgs = MoveArgs{
	Board:      Board{},
	Cubeful:    false,
	ScoreMoves: false,
}

func GetMoves(domain string, moveargs MoveArgs) ([]Move, error) {

	postbody, err := json.Marshal(moveargs)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(domain, "application/json", bytes.NewReader(postbody))
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
