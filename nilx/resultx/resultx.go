package resultx

type Result[T any, E any] struct {
	data  T
	error E
	ok    bool
}

func Ok[T any, E any](data T) Result[T, E] {
	return Result[T, E]{data: data, ok: true}
}

func Err[T any, E any](error E) Result[T, E] {
	return Result[T, E]{error: error, ok: false}
}

func (t Result[T, E]) Match(f func(T) (T, E), g func(E) (T, E)) (T, E) {
	if t.ok {
		return f(t.data)
	}
	return g(t.error)
}
