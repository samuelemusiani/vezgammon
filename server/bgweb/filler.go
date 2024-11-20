package bgweb

import (
	"log/slog"
	"vezgammon/server/types"
)

func fill(m [][]types.Move, g *types.Game) (r [][]types.Move) {
	for i := range m {
		r = append(r, make([]types.Move, len(m[i])))
		copy(r[i], m[i])
	}

	var currentBoard [25]int8

	switch g.CurrentPlayer {
	case types.GameCurrentPlayerP1:
		currentBoard = g.P1Checkers
	case types.GameCurrentPlayerP2:
		currentBoard = g.P2Checkers
	}

	for i := range m {
		for j := 0; j < len(m[i])-1; j++ {
			// A man is outised the board, can't invert moves
			if m[i][j].From == 0 && m[i][j+1].From != 0 {
				continue
			}

			to := m[i][j].To
			from := m[i][j].From
			if to == from && currentBoard[to] == 0 {
				continue
			}

			if m[i][j] == m[i][j+1] {
				continue
			}

			slog.With("i", i, "j", j, "m[i]", m[i]).Debug("Filling")

			tmp := make([]types.Move, len(m[i]))
			copy(tmp, m[i])
			tmp[j], tmp[j+1] = tmp[j+1], tmp[j]

			r = append(r, tmp)
		}
	}

	return
}
