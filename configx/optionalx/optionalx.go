package optionalx

func Parameter[T, R any](constructor func(*T) *R, funcs ...func(*T) *T) *R {
	var p *T
	for _, f := range funcs {
		p = f(p)
	}
	return constructor(p)
}
