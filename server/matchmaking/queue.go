package matchmaking

import (
	"errors"
	"sync"
)

type qel struct {
	User_id int64
	Elo     int64
}

const qlen = 256

var (
	ErrQueueFull       = errors.New("queue is full")
	ErrEmptyQueue      = errors.New("queue is empty")
	ErrElementNotFound = errors.New("element not found")

	queue [qlen]qel

	// point to the first element if start == end it means empty queue
	start int = 0

	// point to the first empty space
	end int = 0

	mutex sync.Mutex
)

func push(e qel) error {
	mutex.Lock()
	if end+1%qlen == start {
		mutex.Unlock()
		return ErrQueueFull
	}

	queue[end] = e
	end = (end + 1) % qlen

	mutex.Unlock()
	return nil
}

func pop() (qel, error) {
	mutex.Lock()
	if isEmpty() {
		mutex.Unlock()
		return qel{}, ErrEmptyQueue
	}

	e := queue[start]
	start = (start + 1) % qlen

	mutex.Unlock()
	return e, nil
}

func remove(e qel) error {
	mutex.Lock()
	if isEmpty() {
		mutex.Unlock()
		return ErrEmptyQueue
	}
	// search element to remove
	for i := range (end + qlen - start) % qlen {
		// element founded
		pos := (i + start) % qlen
		if e == queue[pos] {
			for j := range (end - pos + qlen - 1) % qlen {
				queue[(pos+j)%qlen] = queue[(pos+j+1)%qlen]
			}
			end = (end - 1 + qlen) % qlen
			mutex.Unlock()
			return nil
		}
	}
	mutex.Unlock()
	return ErrElementNotFound
}

func isEmpty() bool {
	return start == end
}
