package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBrewTime(t *testing.T) {
	b := NewBrewer(1)

	start := time.Now()
	c := b.Brew(1, Beans{})
	brewTime := time.Since(start)

	assert.NotNil(t, c)
	assert.Equal(t, 1, c.sizeOunces)
	// brew time should be 1 second-ish
	assert.InDelta(t, brewTime.Milliseconds(), 1, 5)
}
