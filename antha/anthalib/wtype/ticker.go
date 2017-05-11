package wtype

type Ticker struct {
	TickEvery int
	TickBy    int
	Val       int
	tick      int
}

func (t *Ticker) Dup() *Ticker {
	return &Ticker{TickEvery: t.TickEvery, TickBy: t.TickBy, Val: t.Val, tick: t.tick}
}

func (t *Ticker) Tick() int {
	t.tick += 1

	if (t.tick)%t.TickEvery == 0 {
		t.Val += t.TickBy
	}

	return t.Val
}
