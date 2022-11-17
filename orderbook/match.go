package orderbook

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}
