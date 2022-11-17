package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(2_000, sellOrderB)

	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrderA)

	buyOrderA := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrderA)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrderA)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrderA.IsFilled(), true)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	// create three limits on the orderbook
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, ob.BidTotalVolume(), 23.00)

	sellOrder := NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 3.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.bids), 1)

	fmt.Printf("%+v", matches)

}
