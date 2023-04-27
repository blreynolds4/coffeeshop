package models

import (
	"fmt"
	"sync"
)

type OrderChannel chan *Order

type OrderingKiosk interface {
	CreateOrder(name string, item MenuItem) *Order
	setValidity(valid bool)
}

type orderingKiosk struct {
	// the kiosk is only usable for an order if it's valid
	// validity is true to start and when gotten from a pool
	// invalid when put back into the pool
	valid  bool
	orders OrderChannel
}

func newOrderingKiosk(oc OrderChannel) OrderingKiosk {
	return &orderingKiosk{
		valid:  true,
		orders: oc,
	}
}

func (ok *orderingKiosk) setValidity(valid bool) {
	ok.valid = valid
}

func (ok *orderingKiosk) CreateOrder(name string, item MenuItem) *Order {
	if !ok.valid {
		return nil
	}

	// create the order and put it in the shop order channel
	o := NewOrder(name, item)

	fmt.Println(name, "ordered", item.Name)
	ok.orders <- o

	return o
}

type CoffeeShop interface {
	WaitForOrderingKiosk() OrderingKiosk
	LeaveOrderingKiosk(OrderingKiosk)
	Close()
}

type coffeeShop struct {
	Menu      Menu
	baristas  []*barista
	grinders  GrinderPool
	brewers   BrewerPool
	kiosks    KioskPool
	orders    OrderChannel
	closed    bool
	closeWait *sync.WaitGroup
}

func NewCoffeeShop(menu Menu, kioskCount int, baristaCount int, maxBaristaOrders int, grinders GrinderPool, brewers BrewerPool) CoffeeShop {
	result := &coffeeShop{
		Menu:      menu,
		baristas:  make([]*barista, 0, baristaCount),
		grinders:  grinders,
		brewers:   brewers,
		kiosks:    NewKioskPool(),
		orders:    make(OrderChannel, 10*baristaCount),
		closeWait: &sync.WaitGroup{},
	}

	for i := 0; i < kioskCount; i++ {
		result.kiosks.AddKiosk(newOrderingKiosk(result.orders))
	}

	for i := 0; i < baristaCount; i++ {
		name := fmt.Sprintf("Barista-%d", i)
		b := newBarista(name, maxBaristaOrders, result.orders, result.grinders, result.brewers)
		result.baristas = append(result.baristas, b)
		result.closeWait.Add(1)
		go func(b *barista) {
			b.ServeCustomers()
			result.closeWait.Done()
		}(b)
	}

	return result
}

func (cs *coffeeShop) WaitForOrderingKiosk() OrderingKiosk {
	if !cs.closed {
		return cs.kiosks.GetKiosk()
	}

	return nil
}

func (cs *coffeeShop) LeaveOrderingKiosk(k OrderingKiosk) {
	cs.kiosks.AddKiosk(k)
}

func (cs *coffeeShop) Close() {
	// take no more orders
	// anything in progress will finish
	cs.closed = true
	close(cs.orders)
	fmt.Println("Ordering closed")

	// wait for baristas to finish all orders
	fmt.Println("Shop closing, finish up")
	cs.closeWait.Wait()
	fmt.Println("All baristas done")
}
