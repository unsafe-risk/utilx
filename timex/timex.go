package timex

import "time"

func Now() time.Time {
	return time.Now().UTC()
}

func Custom(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).UTC()
}

func TruncateYear(sec int64) time.Time {
	year := Custom(sec, 0).Year()
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TruncateMonth(sec int64) time.Time {
	year := Custom(sec, 0).Year()
	month := Custom(sec, 0).Month()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func TruncateWeek(sec int64) time.Time {
	return Truncate(sec, 1*60*60*24*7)
}

func TruncateDay(sec int64) time.Time {
	return Truncate(sec, 1*60*60*24)
}

func TruncateHour(sec int64) time.Time {
	return Truncate(sec, 1*60*60)
}

func TruncateMin(sec int64) time.Time {
	return Truncate(sec, 1*60)
}

func TruncateSec(sec int64) time.Time {
	return Truncate(sec, 1)
}

func Truncate(sec int64, truncateSec int64) time.Time {
	return time.Unix(sec, 0).UTC().Truncate(time.Second * time.Duration(truncateSec))
}
