package sleepx

import (
	"time"

	"github.com/unsafe-risk/utilx/timex"
)

var (
	// wakeup signal
	wakeup chan struct{} = make(chan struct{}, 1)
)

// sleep for a while && wakeup by signal
func SleepFor(d time.Duration) {
	for {
		select {
		case <-time.After(time.Duration(d) * time.Second):
			return
		case <-wakeup:
			return
		}
	}
}

// sleep until a time && wakeup by signal
func SleepUntil(t time.Time) {
	for {
		select {
		case <-time.After(t.Sub(timex.Now())):
			return
		case <-wakeup:
			return
		}
	}
}

// sleep forever && wakeup by signal
func SleepForever() {
	for {
		select {
		case <-wakeup:
			return
		}
	}
}

// signal to wakeup
func Wakeup() {
	wakeup <- struct{}{}
}
