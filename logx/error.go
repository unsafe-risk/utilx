package logx

import "errors"

var errNoWriter = errors.New("no writer")

func IsNoWriterError(err error) bool {
	return errors.Is(err, errNoWriter)
}

var errNotStruct = errors.New("not a struct")

func IsNotStructError(err error) bool {
	return errors.Is(err, errNotStruct)
}
