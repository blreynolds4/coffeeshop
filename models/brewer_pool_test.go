package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBrewerPool(t *testing.T) {
	b1 := NewBrewer(2)
	bp := NewBrewerPool(b1)

	readBrewer := bp.GetBrewer()

	assert.Equal(t, b1, readBrewer)
}

func TestWaitForBrewer(t *testing.T) {
	b1 := NewBrewer(2)
	bp := NewBrewerPool()

	go func() {
		bp.AddBrewer(b1)
	}()

	readBrewer := bp.GetBrewer()

	assert.Equal(t, b1, readBrewer)
}
