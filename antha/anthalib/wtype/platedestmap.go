package wtype

import (
	"fmt"
)

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

func (pdm Platedestmap) Print() {
	for x := 0; x < 96; x++ {
		fmt.Print(x, " : ")

		for y := 0; y < 96; y++ {
			if len(pdm[x][y]) != 0 {
				fmt.Print(y, "--", len(pdm[x][y]), " ")
			}
		}

		fmt.Println()
	}
}
