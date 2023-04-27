package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGrindTime(t *testing.T) {
	// grind 1g per second
	b := NewGrinder(1)

	start := time.Now()
	grounds := b.Grind(Beans{weightGrams: 1})
	grindTime := time.Since(start)

	assert.Equal(t, 1, grounds.weightGrams)
	// Grind time should be 1 second-ish
	assert.InDelta(t, grindTime.Seconds(), 1, 5)
}
