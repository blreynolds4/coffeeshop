package models

import (
	"fmt"
	"time"
)

type Grinder interface {
	Grind(beans Beans) Beans
}

type grinder struct {
	gramsPerSecond int
}

func NewGrinder(gramsPerSecond int) Grinder {
	return &grinder{
		gramsPerSecond: gramsPerSecond,
	}
}

func (g *grinder) Grind(beans Beans) Beans {
	// Wait for the time it would take to grind the beans
	grindSeconds := g.gramsPerSecond * beans.weightGrams
	fmt.Printf("Grinding %d grams for %d Seconds\n", beans.weightGrams, grindSeconds)
	time.Sleep(time.Duration(grindSeconds) * time.Millisecond)
	fmt.Println("Grind Complete")
	return beans
}
