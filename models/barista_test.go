package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Test each step transistion in coffee making
//

// Initial Order to ReadyToGrind
func TestStartOrder(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	assert.Equal(t, bName, barista.Name)

	order := NewOrder("customer", getTestMenuItem())

	barista.startOrder(order)

	// make sure the grinder request was sent and order is right
	oe, sent := <-barista.activeOrders
	assert.True(t, sent)
	assert.Equal(t, order.Customer, oe.GetOrder().Customer)
	assert.Equal(t, order.Item, oe.GetOrder().Item)
	assert.Equal(t, ReadyToGrind, order.Status)
	assert.Equal(t, order.Status, oe.GetOrder().Status)

	// make sure it's the grinder available event
	gae, ok := oe.(GrinderAvailableEvent)
	assert.True(t, ok)
	assert.NotNil(t, gae.GetGrinder())
}

// ReadyToGrind to Grinding
func TestProgressToGrind(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	order := NewOrder("customer", getTestMenuItem())

	barista.progressOrder(
		NewGrinderAvailableEvent(order, barista.grinders.GetGrinder()))

	// make sure the grind finish was sent and order is right
	oe, sent := <-barista.activeOrders
	assert.True(t, sent)
	assert.Equal(t, order.Customer, oe.GetOrder().Customer)
	assert.Equal(t, order.Item, oe.GetOrder().Item)
	assert.Equal(t, Grinding, order.Status)
	assert.Equal(t, order.Status, oe.GetOrder().Status)

	// make sure it's the grinder available event
	_, ok := oe.(GrindCompleteEvent)
	assert.True(t, ok)
}

// Grinding to ReadyToBrew
func TestProgressToGetBrewer(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	order := NewOrder("customer", getTestMenuItem())

	barista.progressOrder(NewGrindCompleteEvent(order, Beans{weightGrams: 5}))

	// make sure the brewer available was sent and order is right
	oe, sent := <-barista.activeOrders
	assert.True(t, sent)
	assert.Equal(t, order.Customer, oe.GetOrder().Customer)
	assert.Equal(t, order.Item, oe.GetOrder().Item)
	assert.Equal(t, ReadyToBrew, order.Status)
	assert.Equal(t, order.Status, oe.GetOrder().Status)

	// make sure it's the brewer available event
	bae, ok := oe.(BrewerAvailableEvent)
	assert.True(t, ok)
	assert.NotNil(t, bae.GetBrewer())
}

// ReadyToBrew to Brewing
func TestProgressToBrewing(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	order := NewOrder("customer", getTestMenuItem())

	barista.progressOrder(
		NewBrewerAvailableEvent(order, barista.brewers.GetBrewer()))

	// make sure the coffee finish was sent and order is right
	oe, sent := <-barista.activeOrders
	assert.True(t, sent)
	assert.Equal(t, order.Customer, oe.GetOrder().Customer)
	assert.Equal(t, order.Item, oe.GetOrder().Item)
	assert.Equal(t, Brewing, order.Status)
	assert.Equal(t, order.Status, oe.GetOrder().Status)

	// make sure it's the coffee brew complete
	cce, ok := oe.(CoffeeCompleteEvent)
	assert.True(t, ok)
	assert.NotNil(t, cce.GetCoffee())
}

// Brewed to Complete
func TestProgressToComplete(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	order := NewOrder("customer", getTestMenuItem())

	expectedCoffee := &Coffee{sizeOunces: order.Item.Size}

	barista.progressOrder(
		NewCoffeeCompleteEvent(order, expectedCoffee))

	// verify order Wait returns the coffee
	freshCoffee := order.Wait()
	assert.Equal(t, expectedCoffee, freshCoffee)
}

// End to End with ServeCustomer
func TestServeCustomerFullOrder(t *testing.T) {
	bName := "test"
	orderChan := make(OrderChannel)
	barista := newBarista(bName, 1, orderChan, getTestGrinders(), getTestBrewers())

	order := NewOrder("customer", getTestMenuItem())

	expectedCoffee := &Coffee{sizeOunces: order.Item.Size}

	// start the barista
	go barista.ServeCustomers()

	// put in the order
	orderChan <- order

	// verify order Wait returns the coffee
	freshCoffee := order.Wait()

	assert.Equal(t, expectedCoffee, freshCoffee)
}
