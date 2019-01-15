package main

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func makePlateForTest() *wtype.Plate {
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates{X: 127.76, Y: 85.48, Z: 43.1}, welltype, 9.0, 9.0, 0.5, 0.5, 0.5)
	return p
}

func makeTipForTest() *wtype.LHTip {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	return wtype.NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp, 44.7)
}

func makeTipboxForTest() *wtype.LHTipbox {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("ul", 250.0, 10.0, shp, wtype.FlatWellBottom, 7.3, 7.3, 51.2, 0.0, "mm")
	tiptype := makeTipForTest()
	tb := wtype.NewLHTipbox(8, 12, wtype.Coordinates{X: 127.76, Y: 85.48, Z: 120.0}, "me", "mytype", tiptype, w, 9.0, 9.0, 0.5, 0.5, 0.0)
	return tb
}

func makeTipwasteForTest() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "TipwasteForTest", "ACME Corp.", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}
