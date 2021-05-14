package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	knownNow       = 1621003682 // 16:48:02, Friday, 14.05.2021 CEST / 2021-05-14 14:48:02 +0000 UTC
	knownNowUTCStr = "2021-05-14 14:48:02 +0000 UTC"
	knownWeekDay   = 5  // Friday
	knownHour      = 16 // see knownHour
)

func TestBasicTime(t *testing.T) {
	assert.Greater(t, Now(), int64(0))
	assert.Greater(t, Nano(), int64(0))
}

func TestTimeIncrement(t *testing.T) {
	inc := 31 // minutes
	incSeconds := int64(inc * 60)
	now := Now()

	future := IncT(now, inc)

	assert.Greater(t, future, now)
	assert.Equal(t, incSeconds, future-now)

	past := IncT(now, -1*inc)

	assert.Less(t, past, now)
	assert.Equal(t, incSeconds, now-past)
}

func TestElapsedTimeSince(t *testing.T) {
	now := time.Now()
	time.Sleep(2 * time.Second)

	dt := ElapsedTimeSince(now) // dt in ms.

	assert.NotEqual(t, int64(0), now)
	assert.GreaterOrEqual(t, dt, int64(2000)) // 2sec == 2000ms
	assert.Less(t, dt, int64(2200))           // +10% margin, check the order of magnitued
}

func TestToUTC(t *testing.T) {
	utcStr := ToUTC(knownNow)

	assert.NotEmpty(t, utcStr)
	assert.Equal(t, utcStr, knownNowUTCStr)
}

func TestToWeekday(t *testing.T) {
	weekDay := ToWeekday(knownNow)

	assert.Equal(t, knownWeekDay, weekDay)
}

func TestToHour(t *testing.T) {

	// this test might fail if the system clock is syned to UTC ...

	hour := ToHour(knownNow)
	utc := ToHourUTC(knownNow)

	assert.Equal(t, knownHour, hour)
	assert.NotEqual(t, hour, utc)
}
