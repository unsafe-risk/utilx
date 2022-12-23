package closex

import (
	"io"
	"sync"
)

var _ io.Closer = (*CloseManager)(nil)

func New() *CloseManager {
	return &CloseManager{}
}

type CloseManager struct {
	lock sync.Mutex
	list []io.Closer
}

func (c *CloseManager) Append(closer io.Closer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.list = append(c.list, closer)
}

func (c *CloseManager) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	errs := make(ErrorList, 0, len(c.list))

	var err error
	for _, closer := range c.list {
		err = closer.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
