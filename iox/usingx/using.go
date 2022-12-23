package usingx

import (
	"io"
)

func Close(closers ...io.Closer) func(func()) []error {
	return func(f func()) (err []error) {
		errs := []error(nil)
		defer func() {
			for _, closer := range closers {
				if e := closer.Close(); e != nil {
					errs = append(errs, e)
				}
			}
			if len(errs) > 0 {
				err = errs
			}
		}()
		f()
		return
	}
}
