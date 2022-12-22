package optionx

type Option[T any] struct {
	data  T
	exist bool
}

func Some[T any](data T) Option[T] {
	return Option[T]{data: data, exist: true}
}

func None[T any]() Option[T] {
	return Option[T]{exist: false}
}

func (t Option[T]) Match(f func(T) (T, error), g func() (T, error)) (T, error) {
	if t.exist {
		return f(t.data)
	}
	return g()
}
