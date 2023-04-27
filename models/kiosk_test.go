package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKiosk(t *testing.T) {
	// need a channel size or it blocks on write
	orders := make(OrderChannel, 1)

	// kiosk defaults to valid, it gets marked invalid
	// when added to a pool, and valid when taken out
	k := newOrderingKiosk(orders)

	expectedOrder := k.CreateOrder("name", MenuItem{Name: "test"})

	actualOrder := <-orders

	assert.Equal(t, expectedOrder.Customer, actualOrder.Customer)
	assert.Equal(t, expectedOrder.Item.Name, actualOrder.Item.Name)
}

func TestInvalidKiosk(t *testing.T) {
	// need a channel size or it blocks on write
	orders := make(OrderChannel, 1)

	// kiosk defaults to valid, set this to invalid
	// normall set to invalid after an order when it's returned
	// to the kiosk pool for the next customer to get an use
	k := newOrderingKiosk(orders)
	k.setValidity(false)

	// make sure it can't make an order when invalid
	nilOrder := k.CreateOrder("name", MenuItem{Name: "test"})

	assert.Nil(t, nilOrder)
}
