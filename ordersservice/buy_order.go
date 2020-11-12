package ordersservice

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// BuyOrder represents a purchase order for a
// specific stock
type BuyOrder interface {
	// GetID returns the unique identifier used to identify this
	// record in persistent storage
	GetID() string

	// GetPrice returns the price the order was purchased at
	// in increments of 1/1000 of a penny
	GetPrice() int

	// GetQuantity returns the quantity of stock purchased
	GetQuantity() int

	// GetTicker returns the stock ticker code
	GetTicker() string

	// GetTimestamp returns the timestamp from when the order was made
	GetTimestamp() time.Time

	// GetUserID returns the user who will be paying the order
	GetUserID() string
}

// NewBuyOrder constructs a new order instance given a user, ticker, and price.
// it will also