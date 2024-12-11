package bgweb

import (
	"gotest.tools/v3/assert"
	"testing"
	"vezgammon/server/types"
)

func TestFill(t *testing.T) {
	g := &types.Game{}
	g.P1Checkers = [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	g.P2Checkers = [25]int8{0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	g.CurrentPlayer = types.GameCurrentPlayerP1

	m := [][]types.Move{
		{
			types.Move{From: 6, To: 7},
			types.Move{From: 6, To: 7},
			types.Move{From: 6, To: 7},
			types.Move{From: 6, To: 7},
		},
		{
			types.Move{From: 8, To: 9},
			types.Move{From: 6, To: 7},
			types.Move{From: 8, To: 9},
			types.Move{From: 6, To: 7},
		},
		{
			types.Move{From: 8, To: 9},
			types.Move{From: 6, To: 7},
			types.Move{From: 13, To: 14},
			types.Move{From: 6, To: 7},
		},
	}

	fill(m, g)

	g.CurrentPlayer = types.GameCurrentPlayerP2
	r, _ := fill(m, g)
	fill(r, g)

	g.CurrentPlayer = types.GameCurrentPlayerP1
	m = [][]types.Move{
		{
			types.Move{From: 6, To: 7},
			types.Move{From: 6, To: 7},
		},
	}

	fill(m, g)
}

func TestCountMan(t *testing.T) {
	m := []types.Move{
		types.Move{From: 13, To: 14},
		types.Move{From: 13, To: 18},
	}

	c := countMan(7, m)
	assert.Equal(t, c, int8(0))
}

func TestError(t *testing.T) {
	fill(nil, &types.Game{CurrentPlayer: "invalid"})
}

func TestContains(t *testing.T) {
	m := [][]types.Move{
		{
			types.Move{From: 6, To: 7},
			types.Move{From: 6, To: 7},
		},
	}

	c := contains(m, m[0])
	assert.Equal(t, c, true)
}
