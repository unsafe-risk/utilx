package closex

import (
	"errors"
	"fmt"
	"strings"
)

type ErrorList []error

func (list ErrorList) Error() string {
	var buf strings.Builder

	for _, err := range list {
		buf.WriteString(fmt.Sprintf("error - %T\n", err))
		buf.WriteString(err.Error())
		buf.WriteRune('\n')
	}

	return buf.String()
}

func AsErrs(err error) ([]error, bool) {
	var list ErrorList
	if errors.As(err, &list) {
		return list, true
	}

	return nil, false
}
