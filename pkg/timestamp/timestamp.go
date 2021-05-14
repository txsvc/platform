package timestamp

import (
	"time"
)

// Now returns the curent time in seconds UTC
func Now() int64 {
	return time.Now().Unix()
}

// Nano returns the curent time in miliseconds UTC
func Nano() int64 {
	return time.Now().UnixNano()
}

// IncT increments a timstamp (in seconds) by m minutes.
func IncT(t int64, m int) int64 {
	return t + (int64)(m*60)
}

// ElapsedTimeSince returns the difference between t and now.
func ElapsedTimeSince(t time.Time) int64 {
	d := time.Since(t)
	return (int64)(d / time.Millisecond)
}

// ToUTC converts a timestamp to UTC timestamp
func ToUTC(t int64) string {
	return time.Unix(t, 0).UTC().String()
}

// ToWeekday retuns the day of the week for the timestamp
func ToWeekday(t int64) int {
	return int(time.Unix(t, 0).Weekday())
}

// ToHour retuns the hour of the day for the timestamp, in local time
func ToHour(t int64) int {
	return time.Unix(t, 0).Hour()
}

// ToHourUTC retuns the hour of the day for the timestamp, UTC time
func ToHourUTC(t int64) int {
	return time.Unix(t, 0).UTC().Hour()
}
