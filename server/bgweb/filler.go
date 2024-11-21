package bgweb

import (
	"vezgammon/server/types"
)

func fill(m [][]types.Move, g *types.Game) (r [][]types.Move, modified bool) {
	modified = false

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

			//slog.With("i", i, "j", j, "m[i]", m[i]).Debug("Filling")

			tmp := make([]types.Move, len(m[i]))
			copy(tmp, m[i])
			tmp[j], tmp[j+1] = tmp[j+1], tmp[j]

			if !contains(r, tmp) {
				modified = true
				r = append(r, tmp)
			}
		}
	}

	return
}

func contains(a [][]types.Move, b []types.Move) bool {
	if len(b) == 2 {
		for i := range a {
			if a[i][0] == b[0] && a[i][1] == b[1] {
				return true
			}
		}
	} else {
		for i := range a {
			if a[i][0] == b[0] && a[i][1] == b[1] && a[i][2] == b[2] && a[i][3] == b[3] {
				return true
			}
		}
	}

	return false
}
