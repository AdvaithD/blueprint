package main

import (
	"fmt"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

// Limits are like buckets that contain orders at a specific level
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

// Orderbook consists of asks and bids, which are both sorted by price
type Orderbook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %0.2f]", o.Size)
}

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %0.2f || volume: %0.2f]", l.Price, l.TotalVolume)
}
