package orderbook

import (
	"fmt"
	"sort"
)

type Orderbook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceMarketOrder(o *Order) []Match {
	// create a slice of matches
	matches := []Match{}

	// if order is a bid (market buy)
	if o.Bid {
		// make sure size is smaller than total available asks
		if o.Size > ob.AskTotalVolume() {
			// not enough volume to fill
			panic(fmt.Errorf("not enough volume [size: %.2f] to fill [size: %0.2f]", ob.AskTotalVolume(), o.Size))
		}
		// iterate through sorted asks
		for _, limit := range ob.Asks() {
			// fill the current limit with the order 'o'
			limitMatches := limit.Fill(o)
			// append the match to the slice of matches
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(true, limit)
			}
		}
	} else {
		// order is an ask (market sell)
		if o.Size > ob.BidTotalVolume() {
			panic(fmt.Errorf("not enough volume [size: %.2f] to fill [size: %0.2f]", ob.BidTotalVolume(), o.Size))
		}

		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(false, limit)
			}
		}
	}

	return matches
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) []Match {
	// 1. Check if the price level exists
	// 2. If it does, add the order to the price level
	// 3. If it doesn't, create a new price level and add the order to it
	// 4. Add the price level to the orderbook

	// if order is a bid
	if o.Bid {
		if _, ok := ob.BidLimits[price]; !ok {
			l := NewLimit(price)
			ob.BidLimits[price] = l
			ob.bids = append(ob.bids, l)
		}
		ob.BidLimits[price].AddOrder(o)
	} else { // else implies its an ask
		if _, ok := ob.AskLimits[price]; !ok {
			l := NewLimit(price)
			ob.AskLimits[price] = l
			ob.asks = append(ob.asks, l)
		}
		ob.AskLimits[price].AddOrder(o)
	}

	return []Match{}
}

// func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
// 	// 1. Try to match the order with the opposite side of the book
// 	// TODO: Matching logic

// 	// 2. Add rest of order to the books i.e: o.Size is not 0
// 	if o.Size > 0.0 { // meaning its not completely filled
// 		ob.add(price, o)
// 	}

// 	return []Match{}
// }

// return sorted asks
func (ob *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

// return sorted bids
func (ob *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}

func (ob *Orderbook) clearLimit(bid bool, l *Limit) {
	if bid {
		delete(ob.BidLimits, l.Price)
		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l { // if we arrive at the specified limit
				// remove it, slice op, shift oine to the left
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l { // if we arrive at the specified limit
				// remove it
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

// aggregate bid volume on the book
func (ob *Orderbook) BidTotalVolume() float64 {
	total := 0.0

	for i := 0; i < len(ob.bids); i++ {
		total += ob.bids[i].TotalVolume
	}
	return total
}

// agregate ask volume on the book
func (ob *Orderbook) AskTotalVolume() float64 {
	total := 0.0

	for i := 0; i < len(ob.asks); i++ {
		total += ob.asks[i].TotalVolume
	}

	return total
}
