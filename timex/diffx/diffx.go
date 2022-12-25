package diffx

import (
	"time"

	"github.com/unsafe-risk/utilx/timex"
)

type Differ struct {
	start time.Time
	end   time.Time
}

func New() *Differ {
	return &Differ{}
}

func (d *Differ) Start() {
	d.start = timex.Now()
}

func (d *Differ) End() {
	d.end = timex.Now()
}

func (d *Differ) GetDiff() time.Duration {
	return d.end.Sub(d.start)
}

func (d *Differ) GetStart() time.Time {
	return d.start
}

func (d *Differ) GetEnd() time.Time {
	return d.end
}
