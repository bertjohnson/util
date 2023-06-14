package hash

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ComplexObject is a complex object used for testing.
type complexObject struct {
	IgnoredField string `hash:"ignore"`
	StartDate    *time.Time
	Title        string
	Version      int64
}

var (
	now, err = time.Parse(time.RFC1123, "Sun, 29 Jul 2018 10:34:00 CST")
	obj1     = complexObject{
		IgnoredField: "ABC",
		StartDate:    &now,
		Title:        "Object 1",
		Version:      3,
	}
	obj2 = complexObject{
		IgnoredField: "DEF",
		StartDate:    &now,
		Title:        "Object 2",
		Version:      5,
	}
	obj3 = complexObject{
		IgnoredField: "GHI",
		StartDate:    &now,
		Title:        "Object 3",
		Version:      9,
	}
)

// TestHighWayHash tests HighwayHash().
func TestHighwayHash(t *testing.T) {
	assert.NoError(t, err)

	var hash1, hash2, hash3 string
	hash1, err = HighwayHash(context.Background(), obj1)
	assert.NoError(t, err)
	assert.Equal(t, "8u1CYuuRS93", hash1)

	hash2, err = HighwayHash(context.Background(), obj2)
	assert.NoError(t, err)
	assert.Equal(t, "Dh98as2xKUm", hash2)
	assert.NotEqual(t, hash1, hash2)

	hash3, err = HighwayHash(context.Background(), obj3)
	assert.NoError(t, err)
	assert.Equal(t, "6GgqL9yoPvm", hash3)
	assert.NotEqual(t, hash2, hash3)
}

// TestHighwayHashUInt64 tests HighwayHashUInt64().
func TestHighwayHashUInt64(t *testing.T) {
	assert.NoError(t, err)

	var hash1, hash2, hash3 uint64
	hash1, err = HighwayHashUInt64(context.Background(), obj1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x1603dba9d8f3352f), hash1)

	hash2, err = HighwayHashUInt64(context.Background(), obj2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x7a37ddce6285e14b), hash2)
	assert.NotEqual(t, hash1, hash2)

	hash3, err = HighwayHashUInt64(context.Background(), obj3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0xb6588eb07391821f), hash3)
	assert.NotEqual(t, hash2, hash3)
}
