package wtype

type Platedestmap [][][]*LHInstruction

func NewPlatedestmap() Platedestmap {
	h := make([][][]*LHInstruction, 96)

	for x := 0; x < 96; x++ {
		h[x] = make([][]*LHInstruction, 96)

		for y := 0; y < 96; y++ {
			h[x][y] = make([]*LHInstruction, 0, 2)
		}
	}

	return h
}
