package diffx

import (
	"fmt"
	"testing"

	"github.com/unsafe-risk/utilx/timex/sleepx"
)

func TestDiff(t *testing.T) {
	differ := New()

	differ.Start()

	// logic to measure elapsed time
	sleepx.New().SleepFor(1)

	differ.End()

	fmt.Printf("diff : %s\n", differ.GetDiff())
}
