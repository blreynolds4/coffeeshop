# coffeeshop

Model the workings of a coffee shop.

Implementation Notes:

1. The model assumes the shop has enough coffee to fill all the orders.
2. There is a maximum number of orders a barista can handle.  At that point they focus on what they have.
3. Working in real seconds made runs take a very long time.  I kept the seconds labels but internally use milliseconds for grinding and brewing.

## Building

go build -o coffee-sim

## Running

./coffee-sim  --help
Usage of ./coffee-sim:
  -barista-count int
        The count of baristas working in the coffee shop (default 1)
  -barista-order-count int
        The maximum number of orders a barista can work on at a time (default 5)
  -brewer-count int
        The count of brewers in the coffee shop (default 1)
  -customer-count int
        The count of customers ordering in the coffee shop (default 1)
  -grinder-count int
        The count of grinders in the coffee shop (default 1)
  -kiosk-count int
        The count of ordering kiosks in the coffee shop (default 1)

Example:
  coffee-sim -barista-count 2 -barista-order-count 10 -brewer-count 3 -grinder-count 3 -kiosk-count 2 -customer-count 20
