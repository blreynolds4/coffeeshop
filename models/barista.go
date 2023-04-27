package models

import (
	"fmt"
	"sync"
)

type OrderStepsChannel chan OrderEvent

// Barista represents the process of making coffee:
// getting a grinder
// doing the grind
// getting a brewer
// doing the brew
// notifying the customer
// A barista can work on many coffees at a time while
// waiting on each step.  If they work on more than maxActiveOrders
// they will not be able to take any new orders
type barista struct {
	Name         string
	newOrders    OrderChannel
	activeOrders OrderStepsChannel
	grinders     GrinderPool
	brewers      BrewerPool
	orderCount   int
	countLock    *sync.Mutex
}

func newBarista(name string, maxActiveOrders int, newOrders OrderChannel, g GrinderPool, b BrewerPool) *barista {
	return &barista{
		Name:         name,
		newOrders:    newOrders,
		activeOrders: make(OrderStepsChannel, maxActiveOrders),
		grinders:     g,
		brewers:      b,
		orderCount:   0,
		countLock:    &sync.Mutex{},
	}
}

// Keep an internal count of active orders to help
// the barista know when it's done
func (b *barista) incrementOrderCount() {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	b.orderCount += 1
}

func (b *barista) decrementOrderCount() {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	b.orderCount -= 1
}

func (b *barista) getCurrentOrderCount() int {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	return b.orderCount
}

// ServeCustomers reads the new orders channel to start new orders
// or reads current orders channel to progress existing orders.
// If a stop is requested (by closing the shop and new orders channel)
// focus on the existing orders till they're done.
func (b *barista) ServeCustomers() {
	for stopping := false; !stopping || b.getCurrentOrderCount() > 0; {
		if !stopping {
			select {
			case existingOrderEvent := <-b.activeOrders:
				b.progressOrder(existingOrderEvent)
			case newOrder, isOpen := <-b.newOrders:
				if newOrder == nil && !isOpen {
					stopping = true
				} else {
					b.startOrder(newOrder)
				}
			}
		} else {
			// we're stopping, only read from existing orders
			// and only if order count is > 0
			if b.getCurrentOrderCount() > 0 {
				existingOrderEvent := <-b.activeOrders
				b.progressOrder(existingOrderEvent)
			}
		}
	}

	fmt.Println(b.Name, "is done for the day")
}

func (b *barista) startOrder(newOrder *Order) {
	fmt.Println(b.Name, "is working on order from", newOrder.Customer)
	// new orders need to be ground, set the status to ReadyToGrind
	// request a grinder and move on till it's available
	newOrder.Status = ReadyToGrind
	b.incrementOrderCount()
	go func() {
		grinder := b.grinders.GetGrinder()
		fmt.Println(b.Name, "got grinder for", newOrder.Customer)
		// notfiy the barista the grinder is available
		b.activeOrders <- NewGrinderAvailableEvent(newOrder, grinder)
	}()
}

func (b *barista) progressOrder(event OrderEvent) {
	switch {
	case isGrinderAvailable(event):
		b.grindCoffee(event.(GrinderAvailableEvent))

	case isGrindComplete(event):
		b.requestBrewer(event.(GrindCompleteEvent))

	case isBrewerAvailable(event):
		b.brewCoffee(event.(BrewerAvailableEvent))

	case isCoffeeComplete(event):
		b.serveCoffee(event.(CoffeeCompleteEvent))
	}
}

func (b *barista) grindCoffee(ge GrinderAvailableEvent) {
	go func() {
		order := ge.GetOrder()
		order.Status = Grinding
		grinder := ge.GetGrinder()

		// grind the right amount of beans for the order
		ungroundBeans := Beans{weightGrams: order.Item.CoffeeRatio * order.Item.Size}

		// save the ground beans
		fmt.Println(b.Name, "is grinding coffee for", order.Customer)
		beans := grinder.Grind(ungroundBeans)
		b.grinders.AddGrinder(grinder)

		b.activeOrders <- NewGrindCompleteEvent(order, beans)
	}()
}

func (b *barista) requestBrewer(ge GrindCompleteEvent) {
	order := ge.GetOrder()
	order.Status = ReadyToBrew
	order.GroundBeans = ge.GetBeans()
	fmt.Println(b.Name, "is getting a brewer for", order.Customer)
	go func() {
		brewer := b.brewers.GetBrewer()
		fmt.Println(b.Name, "got a brewer for", order.Customer)
		b.activeOrders <- NewBrewerAvailableEvent(order, brewer)
	}()
}

func (b *barista) brewCoffee(ge BrewerAvailableEvent) {
	go func() {
		order := ge.GetOrder()
		order.Status = Brewing
		brewer := ge.GetBrewer()
		fmt.Println(b.Name, "is brewing coffee for", order.Customer)

		// brew the coffee to the final volume
		coffee := brewer.Brew(order.Item.Size, order.GroundBeans)
		fmt.Println(b.Name, "is done brewing coffee for", order.Customer)
		b.brewers.AddBrewer(brewer)

		b.activeOrders <- NewCoffeeCompleteEvent(order, coffee)
	}()
}

func (b *barista) serveCoffee(cc CoffeeCompleteEvent) {
	b.decrementOrderCount()
	order := cc.GetOrder()
	order.Status = Complete
	coffee := cc.GetCoffee()

	fmt.Println(b.Name, "says coffee is ready for", order.Customer)
	order.NotifyCustomer(coffee)
}
