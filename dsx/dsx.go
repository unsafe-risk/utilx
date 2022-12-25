package dsx

type Deque[T any] interface {
	PushFront(T)
	PushBack(T)
	PopFront() (T, bool)
	PopBack() (T, bool)
	PeekFront() (T, bool)
	PeekBack() (T, bool)
}

type List[T any] interface {
	Append(T)
	Insert(int, T)
	Remove(int) (T, bool)
	At(int) (T, bool)
	Iterate(func(T) bool)
	Len() int
}
