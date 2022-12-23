package closex

import "io"

var Default = New()

func Append(closer io.Closer) {
	Default.Append(closer)
}

func Close() error {
	return Default.Close()
}
