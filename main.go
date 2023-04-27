package main

import (
	"blreynolds4/coffeeshop/models"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	RegularSizeOunces = 8
	LargeSizeOunces   = 12

	// common ratio is 1 gram coffee per 17ml of water
	// 17 ml is .6 ounces, scale that to 1 ounce
	// is roughly 28ml so we'll call it 2 grams per ounce
	RegularBrewingRatioGramsPerOunce = 2
	StrongBrewingRatioGramsPerOunce  = 4
)

func main() {
	var cliGrinderCount int
	var cliBrewerCount int
	var cliKioskCount int
	var cliBaristaCount int
	var cliCustomerCount int
	var cliBaristaOrderCount int

	flag.IntVar(&cliGrinderCount, "grinder-count", 1, "The count of grinders in the coffee shop")
	flag.IntVar(&cliBrewerCount, "brewer-count", 1, "The count of brewers in the coffee shop")
	flag.IntVar(&cliKioskCount, "kiosk-count", 1, "The count of ordering kiosks in the coffee shop")
	flag.IntVar(&cliBaristaCount, "barista-count", 1, "The count of baristas working in the coffee shop")
	flag.IntVar(&cliCustomerCount, "customer-count", 1, "The count of customers ordering in the coffee shop")
	flag.IntVar(&cliBaristaOrderCount, "barista-order-count", 5, "The maximum number of orders a barista can work on at a time")

	// parse command line
	flag.Parse()

	// create a menu for the coffee shop
	// this could be proivded via a config file
	menu := models.Menu{
		models.MenuItem{
			Name:        "Regular",
			Size:        RegularSizeOunces,
			CoffeeRatio: RegularBrewingRatioGramsPerOunce,
		},
		models.MenuItem{
			Name:        "Regular Strong",
			Size:        RegularSizeOunces,
			CoffeeRatio: StrongBrewingRatioGramsPerOunce,
		},
		models.MenuItem{
			Name:        "Large Regular",
			Size:        LargeSizeOunces,
			CoffeeRatio: RegularBrewingRatioGramsPerOunce,
		},
		models.MenuItem{
			Name:        "Large Strong",
			Size:        LargeSizeOunces,
			CoffeeRatio: StrongBrewingRatioGramsPerOunce,
		},
	}

	// Create pool of grinders.  They grind in grams per second
	grinders := models.NewGrinderPool()
	for i := 0; i < cliGrinderCount; i++ {
		// create a grinder with up to 10 grams per second speed
		grinders.AddGrinder(models.NewGrinder(rand.Intn(10)))
	}

	// Create pool of brewers.  They brew in ounces per second
	brewers := models.NewBrewerPool()
	for i := 0; i < cliBrewerCount; i++ {
		// create brewer with up to LargeSizeOunces per second
		brewers.AddBrewer(models.NewBrewer(rand.Intn(LargeSizeOunces)))
	}

	// create the coffee shop with all the stuff
	shop := models.NewCoffeeShop(menu, cliKioskCount, cliBaristaCount, cliBaristaOrderCount, grinders, brewers)

	orderWaitGroup := sync.WaitGroup{}
	orderWaitGroup.Add(cliCustomerCount)
	start := time.Now()
	for i := 0; i < cliCustomerCount; i++ {
		// in parallel, all at once, make calls to MakeCoffee
		go func(customer string) {
			// model the customer
			// wait for turn to order
			// order random coffee off menu
			// leave the kiosk for the next person
			kiosk := shop.WaitForOrderingKiosk()
			order := kiosk.CreateOrder(customer, menu[rand.Intn(len(menu))])
			shop.LeaveOrderingKiosk(kiosk)

			order.Wait()
			fmt.Println(customer + " says Thank You")
			orderWaitGroup.Done()
		}(fmt.Sprintf("Customer-%d", i))
	}

	fmt.Println("Waiting for all customers to order...")

	// wait for orders to be put in before we close the
	// shop
	orderWaitGroup.Wait()
	fmt.Println("Customers have all ordered.")

	// stop taking orders and wait for baristas to finish
	shop.Close()
	fmt.Println("All orders complete.")
	runTime := time.Since(start)
	fmt.Println("Run time", runTime)
	fmt.Println("Avg Coffee time", runTime.Milliseconds()/int64(cliCustomerCount))
}

// Premise: we want to model a coffee shop. An order comes in, and then with a limited amount of grinders and
// brewers (each of which can be "busy"): we must grind unground beans, take the resulting ground beans, and then
// brew them into liquid coffee. We need to coordinate the work when grinders and/or brewers are busy doing work
// already. What Go datastructure(s) might help us coordinate the steps: order -> grinder -> brewer -> coffee?
//
// Some of the struct types and their functions need to be filled in properly. It may be helpful to finish the
// Grinder impl, and then Brewer impl each, and then see how things all fit together inside CoffeeShop afterwards.

// Issues with the above
// 1. Assumes that we have unlimited amounts of grinders and brewers.
//		- How do we build in logic that takes into account that a given Grinder or Brewer is busy?
// 2. Does not take into account that brewers must be used after grinders are done.
// 		- Making a coffee needs to be done sequentially: find an open grinder, grind the beans, find an open brewer,
//		  brew the ground beans into coffee.
// 3. A lot of assumptions (i.e. 2 grams needed for 1 ounce of coffee) are left as comments in the code.
// 		- How can we make these assumptions configurable, so that our coffee shop can serve let's say different
//		  strengths of coffee via the Order that is placed (i.e. 5 grams of beans to make 1 ounce of coffee)?
