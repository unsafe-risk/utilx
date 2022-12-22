package timex

import "time"

func Now() time.Time {
	return time.Now().UTC()
}

func Custom(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).UTC()
}
