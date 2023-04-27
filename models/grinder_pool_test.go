package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGrinderPool(t *testing.T) {
	g1 := NewGrinder(4)
	gp := NewGrinderPool(g1)

	readGrinder := gp.GetGrinder()
	assert.NotNil(t, readGrinder)
	assert.Equal(t, g1, readGrinder)
}

func TestWaitForGrinder(t *testing.T) {
	g1 := NewGrinder(4)
	gp := NewGrinderPool()

	go func() {
		gp.AddGrinder(g1)
	}()

	readGrinder := gp.GetGrinder()

	assert.NotNil(t, readGrinder)
	assert.Equal(t, g1, readGrinder)
}
