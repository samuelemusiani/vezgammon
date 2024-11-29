package matchmaking

import (
	"errors"
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
)

func push(e qel) error {

	if end+1%qlen == start {

		return ErrQueueFull
	}

	queue[end] = e
	end = (end + 1) % qlen

	return nil
}

func pop() (qel, error) {

	if isEmpty() {

		return qel{}, ErrEmptyQueue
	}

	e := queue[start]
	start = (start + 1) % qlen

	return e, nil
}

func remove(e qel) error {

	if isEmpty() {

		return ErrEmptyQueue
	}
	// search element to remove
	for i := range length() {
		// element founded
		pos := (i + start) % qlen
		if e.User_id == queue[pos].User_id {
			for j := range (end - pos + qlen - 1) % qlen {
				queue[(pos+j)%qlen] = queue[(pos+j+1)%qlen]
			}
			end = (end - 1 + qlen) % qlen

			return nil
		}
	}

	return ErrElementNotFound
}

func isEmpty() bool {
	return start == end
}

func length() int {
	return (end - start + qlen) % qlen
}
