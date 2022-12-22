package ringx

import "fmt"

func NewRing(cap int64) *Ring {
	if cap == 0 {
		return nil
	}
	return &Ring{
		data: make([]interface{}, cap),
		cap:  cap,
		head: 0,
		tail: 0,
	}
}

type Ring struct {
	data []interface{}
	cap  int64
	head int64
	tail int64
}

func (t *Ring) Cap() int64 {
	return t.cap
}

func (t *Ring) IsEmpty() bool {
	if t.head == t.tail {
		return true
	}
	return false
}

func (t *Ring) IsFull() bool {
	if t.head == (t.tail+1)%t.cap {
		return true
	}
	return false
}

func (t *Ring) Enqueue(data interface{}) error {
	if t.IsFull() == true {
		return fmt.Errorf("Queue is full | cap - %d", t.cap)
	}
	t.data[t.tail] = data
	t.tail = (t.tail + 1) % t.cap
	return nil
}

func (t *Ring) Dequeue() (data interface{}) {
	if t.IsEmpty() {
		return nil
	}
	data = t.data[t.head]
	t.head = (t.head + 1) % t.cap
	return data
}

func (t *Ring) Head() interface{} {
	return t.data[t.head]
}
