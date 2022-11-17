package orderbook

import (
	"fmt"
	"time"
)

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o *Order) String() string {
	return fmt.Sprintf("[size: %0.2f]", o.Size)
}

func (o *Order) IsFilled() bool {
	return o.Size == 0.0
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}
