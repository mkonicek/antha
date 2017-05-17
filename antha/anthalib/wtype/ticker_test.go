package wtype

import "testing"

func TestTicker(t *testing.T) {
	testTicker(t, 1, 1, 0)
	testTicker(t, 1, 2, 0)
	testTicker(t, 2, 1, 0)
	testTicker(t, 1, 1, 3)
	testTicker(t, 1, 2, 3)
	testTicker(t, 2, 1, 3)
}

func testTicker(t *testing.T, te, tb, v int) {

	ticker := Ticker{TickEvery: te, TickBy: tb, Val: v}

	expected := v
	lastincr := 0

	for i := 0; i < 100; i++ {
		if ticker.Val != expected {
			t.Errorf("TickEvery: %d TickBy: %d Expected %d got %d", ticker.TickEvery, ticker.TickBy, i, ticker.Val)
		}

		ticker.Tick()

		lastincr += 1

		if lastincr == te {
			expected += tb
			lastincr = 0
		}
	}
}
