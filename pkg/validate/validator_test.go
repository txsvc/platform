package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := New("test")
	assert.NotNil(t, v)
	assert.Equal(t, "test", v.Name)
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestStringEquals(t *testing.T) {
	v := New("test")

	v.StringEquals("a", "a")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.StringEquals("a", "b")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestStringNotEmpty(t *testing.T) {
	v := New("test")

	v.StringNotEmpty("a", "a")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.StringNotEmpty("", "b")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestNotNil(t *testing.T) {
	v1 := New("test")
	v1.NotNil("a")
	assert.Equal(t, 0, v1.Errors)
	assert.Equal(t, 0, v1.Warnings)

	v2 := New("test")
	v2.NotNil("")
	assert.Equal(t, 0, v2.Errors)
	assert.Equal(t, 0, v2.Warnings)

	var s1 *Assertion
	v3 := New("test")
	v3.NotNil(s1)
	assert.Equal(t, 1, v3.Errors)
	assert.Equal(t, 0, v3.Warnings)

	v4 := New("test")
	s2 := new(Assertion)
	v4.NotNil(s2)
	assert.Equal(t, 0, v4.Errors)
	assert.Equal(t, 0, v4.Warnings)
}

func TestNonZero(t *testing.T) {
	v := New("test")

	v.NonZero(42, "a")
	assert.Equal(t, 0, v.Errors)
	v.NonZero(-66, "a")
	assert.Equal(t, 0, v.Errors)
	v.NonZero(0, "a")
	assert.Equal(t, 1, v.Errors)
}

func TestISO639(t *testing.T) {
	v := New("test")

	v.ISO639("en_US")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)
	v.ISO639("en_us")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestRFC1123Z(t *testing.T) {
	v := New("test")

	v.RFC1123Z("Fri, 04 Jun 2021 06:57:26 +0000")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.RFC1123Z("04 Jun 2021 06:57:26 +0000")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.RFC1123Z("Fri, 04 Jun 2021 06:57:26")
	assert.Equal(t, 2, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestTimestamp(t *testing.T) {
	v := New("test")

	v.Timestamp("1")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.Timestamp("-1")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.Timestamp("sdfdfd")
	assert.Equal(t, 2, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestMapNotEmpty(t *testing.T) {

	m := make(map[string]string)

	v1 := New("test")
	v1.MapNotEmpty(m, "a")
	assert.Equal(t, 1, v1.Errors)
	assert.Equal(t, 0, v1.Warnings)

	v2 := New("test")
	m["foo"] = "bar"
	v2.MapNotEmpty(m, "a")
	assert.Equal(t, 0, v2.Errors)
	assert.Equal(t, 0, v2.Warnings)
}

func TestMapContains(t *testing.T) {
	v := New("test")
	m := make(map[string]string)
	m["foo"] = "bar"

	v.MapContains(m, "foo", "a")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.MapContains(m, "bar", "a")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}
