package matchmaking

import (
	"sync"
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

func TestRemove(t *testing.T) {
	// to empty the queue
	start = end

	e := qel{User_id: 69, Elo: 1000}
	err := push(e)
	assert.NilError(t, err)

	err = remove(e)
	assert.NilError(t, err)

	_, err = pop()
	assert.ErrorIs(t, err, ErrEmptyQueue)
}

func TestEveryElement(t *testing.T) {
	var a [qlen - 1]qel

	for i := range qlen - 1 {
		a[i] = qel{User_id: rand.Int64N(1000), Elo: rand.Int64N(2000)}
		err := push(a[i])
		assert.NilError(t, err)
	}

	for i := range qlen - 1 {
		e, err := pop()
		assert.NilError(t, err)
		assert.DeepEqual(t, a[i], e)
	}
}

func TestPartialQueueFilling(t *testing.T) {
	start = 0
	end = 0

	elementsToFill := qlen / 2
	for i := 0; i < elementsToFill; i++ {
		e := qel{User_id: int64(i), Elo: int64(1000 + i)}
		err := push(e)
		assert.NilError(t, err)
	}

	assert.Equal(t, end, elementsToFill)
}

func TestRemoveNonexistentElement(t *testing.T) {
	start = 0
	end = 0

	e := qel{User_id: 123, Elo: 1500}
	err := remove(e)
	assert.ErrorIs(t, err, ErrEmptyQueue)

	err = push(qel{User_id: 456, Elo: 1600})
	assert.NilError(t, err)

	err = remove(e)
	assert.ErrorIs(t, err, ErrElementNotFound)
}

func TestConcurrentOperations(t *testing.T) {
	start = 0
	end = 0

	numOps := 100
	var wg sync.WaitGroup
	wg.Add(numOps * 2)

	for i := 0; i < numOps; i++ {
		go func(id int) {
			defer wg.Done()
			err := push(qel{User_id: int64(id), Elo: int64(1000 + id)})
			if err != nil && err != ErrQueueFull {
				t.Errorf("Unexpected error: %v", err)
			}
		}(i)
	}

	for range numOps {
		go func() {
			defer wg.Done()
			_, err := pop()
			if err != nil && err != ErrEmptyQueue {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()
}
