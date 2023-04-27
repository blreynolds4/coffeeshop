package models

import "sync"

type Beans struct {
	weightGrams int
}

type Coffee struct {
	sizeOunces int
}

type MenuItem struct {
	Name        string
	Size        int
	CoffeeRatio int
}

type Menu []MenuItem

type OrderStatus int

const (
	Ordered OrderStatus = iota
	ReadyToGrind
	Grinding
	ReadyToBrew
	Brewing
	Complete
)

func (s OrderStatus) String() string {
	switch s {
	case Ordered:
		return "Ordered"
	case ReadyToGrind:
		return "Ready to Grind"
	case Grinding:
		return "Grinding"
	case ReadyToBrew:
		return "Ready to Brew"
	case Brewing:
		return "Brewing"
	case Complete:
		return "Complete"
	}

	return "unknown order status"
}

type Order struct {
	Customer    string
	Item        MenuItem
	Status      OrderStatus
	GroundBeans Beans
	freshCoffee *Coffee
	doneFlag    *sync.Cond
}

func NewOrder(cust string, item MenuItem) *Order {
	return &Order{
		Customer: cust,
		Item:     item,
		Status:   Ordered,
		doneFlag: sync.NewCond(&sync.Mutex{}),
	}
}

type OrderEvent interface {
	GetOrder() *Order
}

type GrinderAvailableEvent interface {
	OrderEvent
	GetGrinder() Grinder
}

type GrindCompleteEvent interface {
	OrderEvent
	GetBeans() Beans
}

type BrewerAvailableEvent interface {
	OrderEvent
	GetBrewer() Brewer
}

type CoffeeCompleteEvent interface {
	OrderEvent
	GetCoffee() *Coffee
}

type orderEvent struct {
	order *Order
}

type grinderAvailable struct {
	orderEvent
	grinder Grinder
}

type grindComplete struct {
	orderEvent
	beans Beans
}
type brewerAvailable struct {
	orderEvent
	brewer Brewer
}

type coffeeComplete struct {
	orderEvent
	coffee *Coffee
}

func NewGrinderAvailableEvent(o *Order, g Grinder) GrinderAvailableEvent {
	return &grinderAvailable{
		orderEvent: orderEvent{order: o},
		grinder:    g,
	}
}

func NewGrindCompleteEvent(o *Order, b Beans) GrindCompleteEvent {
	return &grindComplete{
		orderEvent: orderEvent{order: o},
		beans:      b,
	}
}

func NewBrewerAvailableEvent(o *Order, b Brewer) BrewerAvailableEvent {
	return &brewerAvailable{
		orderEvent: orderEvent{order: o},
		brewer:     b,
	}
}

func NewCoffeeCompleteEvent(o *Order, c *Coffee) CoffeeCompleteEvent {
	return &coffeeComplete{
		orderEvent: orderEvent{order: o},
		coffee:     c,
	}
}

func (oe *orderEvent) GetOrder() *Order {
	return oe.order
}

func (ga *grinderAvailable) GetGrinder() Grinder {
	return ga.grinder
}

func (gc *grindComplete) GetBeans() Beans {
	return gc.beans
}

func (ga *brewerAvailable) GetBrewer() Brewer {
	return ga.brewer
}

func (cc *coffeeComplete) GetCoffee() *Coffee {
	return cc.coffee
}

func isGrinderAvailable(e OrderEvent) bool {
	_, isGae := e.(GrinderAvailableEvent)
	return isGae
}

func isGrindComplete(e OrderEvent) bool {
	_, isGae := e.(GrindCompleteEvent)
	return isGae
}

func isBrewerAvailable(e OrderEvent) bool {
	_, isBae := e.(BrewerAvailableEvent)
	return isBae
}

func isCoffeeComplete(e OrderEvent) bool {
	_, isCc := e.(CoffeeCompleteEvent)
	return isCc
}

func (o *Order) NotifyCustomer(c *Coffee) {
	// release the wait with the coffee
	o.doneFlag.L.Lock()
	defer o.doneFlag.L.Unlock()

	o.freshCoffee = c

	o.doneFlag.Broadcast()
}

func (o *Order) Wait() *Coffee {
	o.doneFlag.L.Lock()

	if o.freshCoffee == nil {
		// wait for the coffee  to be Completed
		o.doneFlag.Wait()
	}

	// return the coffee put in by Notify
	return o.freshCoffee
}
