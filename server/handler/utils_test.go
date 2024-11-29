package handler

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestEloCalculation(t *testing.T) {
	elo1, elo2 := calculateElo(1000, 1000, true)
	assert.Equal(t, elo1, int64(1016))
	assert.Equal(t, elo2, int64(984))
}
