package matchmaking

import (
	"testing"

	"math/rand/v2"

	"gotest.tools/v3/assert"
)

func TestEmptyQueue(t *testing.T) {
	assert.Equal(t, isEmpty(), true)
}

func TestPushPop(t *testing.T) {
	e := qel{User_id: 69, Elo: 1000}
	err := push(e)
	assert.NilError(t, err)

	e2, err := pop()
	assert.NilError(t, err)

	assert.DeepEqual(t, e, e2)
}

func TestErrors(t *testing.T) {
	_, err := pop()
	assert.ErrorIs(t, err, ErrEmptyQueue)

	for range qlen - 1 {
		id := rand.Int64N(1000)
		elo := rand.Int64N(2000)
		err := push(qel{id, elo})
		assert.NilError(t, err)

	}
	id := rand.Int64N(1000)
	elo := rand.Int64N(2000)
	err = push(qel{id, elo})
	assert.ErrorIs(t, err, ErrQueueFull)

}
