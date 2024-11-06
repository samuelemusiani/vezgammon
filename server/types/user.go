package types

import (
	"math/rand"
)

func NewDices() [2]int {
	var dices [2]int
	for i := 0; i < len(dices); i++ {
		dices[i] = rand.Intn(6) + 1
	}
	return dices
}
