package models

import (
	"fmt"
	"time"
)

type Brewer interface {
	Brew(finishedVolume int, beans Beans) *Coffee
}

type brewer struct {
	// assume we have unlimited water, but we can only run a
	// certain amount of water per second into our brewer + beans
	ouncesWaterPerSecond int
}

func NewBrewer(ouncesWaterPerSecond int) Brewer {
	return &brewer{
		ouncesWaterPerSecond: ouncesWaterPerSecond,
	}
}

// take the right amount of water and the beans and brew the coffee
func (b *brewer) Brew(finishedVolume int, beans Beans) *Coffee {
	// do the brewing
	brewTime := b.ouncesWaterPerSecond * finishedVolume
	fmt.Printf("Brewing %d grams for %d Seconds\n", beans.weightGrams, brewTime)
	time.Sleep(time.Duration(brewTime) * time.Millisecond)
	fmt.Println("Brew Complete")
	return &Coffee{sizeOunces: finishedVolume}
}
