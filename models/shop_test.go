package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderingAndKioskValidity(t *testing.T) {
	shop := NewCoffeeShop(Menu{getTestMenuItem()},
		1, // ordering kiosk count
		1, // barista count
		1, // max orders per barista
		getTestGrinders(),
		getTestBrewers())

	kiosk := shop.WaitForOrderingKiosk()
	assert.NotNil(t, kiosk)

	// make the order and verify it makes sense
	order := kiosk.CreateOrder("test", getTestMenuItem())
	assert.NotNil(t, order)
	assert.Equal(t, "test", order.Customer)
	assert.Equal(t, getTestMenuItem(), order.Item)

	// leave the kiosk and verify we can't order from it
	shop.LeaveOrderingKiosk(kiosk)
	assert.Nil(t, kiosk.CreateOrder("nope", getTestMenuItem()))

	// close the shop, orders should complete
	shop.Close()

	// wait for the order to be complete and make sure it's right
	c := order.Wait()
	assert.NotNil(t, c)
	assert.Equal(t, getTestMenuItem().Size, c.sizeOunces)

	// verify can't get ordering kiosk after close
	assert.Nil(t, shop.WaitForOrderingKiosk())
}
