package orderbook

import (
	"fmt"
	"sort"
)

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %0.2f || volume: %0.2f]", l.Price, l.TotalVolume)
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) Fill(o *Order) []Match {
	matches := []Match{}

	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)

		if o.IsFilled() {
			break
		}
	}

	return matches
}

// one of the orders is going to be completely eaten up
// the other will be partially filled
func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)

	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}

	if a.Size > b.Size {
		// b is completely eaten up
		a.Size -= b.Size
		sizeFilled = b.Size
		b.Size = 0.0
	} else {
		// a is completely eaten up
		b.Size -= a.Size
		sizeFilled = a.Size
		a.Size = 0.0
	}

	return Match{
		Bid:        bid,
		Ask:        ask,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			// l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	// sort remaining orders by timestamp
	sort.Sort(l.Orders)
}
