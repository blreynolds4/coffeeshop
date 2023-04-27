package models

import "sync"

type sharedPool[A any] struct {
	items  []A
	signal sync.Cond
}

func (sp *sharedPool[A]) AddToPool(obj A) {
	sp.signal.L.Lock()
	defer sp.signal.L.Unlock()

	sp.items = append(sp.items, obj)
	sp.signal.Signal()
}

func (gp *sharedPool[A]) GetFromPool() A {
	gp.signal.L.Lock()
	defer gp.signal.L.Unlock()

	for len(gp.items) == 0 {
		gp.signal.Wait()
	}

	result := gp.items[0]
	gp.items = gp.items[1:]
	return result
}

type GrinderPool interface {
	AddGrinder(Grinder)
	GetGrinder() Grinder
}

type grinderPool struct {
	sharedPool[Grinder]
}

func NewGrinderPool(grinders ...Grinder) GrinderPool {
	result := &grinderPool{
		sharedPool[Grinder]{
			items:  make([]Grinder, 0, len(grinders)),
			signal: *sync.NewCond(&sync.Mutex{}),
		},
	}

	for _, g := range grinders {
		result.AddToPool(g)
	}

	return result
}

func (gp *grinderPool) AddGrinder(g Grinder) {
	gp.AddToPool(g)
}

func (gp *grinderPool) GetGrinder() Grinder {
	return gp.GetFromPool()
}

type BrewerPool interface {
	AddBrewer(Brewer)
	GetBrewer() Brewer
}

type brewerPool struct {
	sharedPool[Brewer]
}

func NewBrewerPool(brewers ...Brewer) BrewerPool {
	result := &brewerPool{
		sharedPool[Brewer]{
			items:  make([]Brewer, 0, len(brewers)),
			signal: *sync.NewCond(&sync.Mutex{}),
		},
	}

	for _, b := range brewers {
		result.AddToPool(b)
	}

	return result
}

func (bp *brewerPool) AddBrewer(b Brewer) {
	bp.AddToPool(b)
}

func (bp *brewerPool) GetBrewer() Brewer {
	return bp.GetFromPool()
}

type kioskPool struct {
	sharedPool[OrderingKiosk]
}

type KioskPool interface {
	AddKiosk(OrderingKiosk)
	GetKiosk() OrderingKiosk
}

func NewKioskPool() KioskPool {
	result := &kioskPool{
		sharedPool[OrderingKiosk]{
			items:  make([]OrderingKiosk, 0),
			signal: *sync.NewCond(&sync.Mutex{}),
		},
	}

	return result
}

func (kp *kioskPool) AddKiosk(ok OrderingKiosk) {
	// mark the kiosk invalid until it is waited for again
	ok.setValidity(false)
	kp.AddToPool(ok)
}

func (kp *kioskPool) GetKiosk() OrderingKiosk {
	// the kiosks from the pool are valid
	result := kp.GetFromPool()
	result.setValidity(true)
	return result
}
