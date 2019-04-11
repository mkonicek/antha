package cache

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/laboratory/testlab"
)

const plateType = "pcrplate"

func TestPlateReuse(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	inv := testlab.InventoryWithSharedPlates()
	cache := plateCache.NewPlateCache(inv.PlateTypes)

	firstPlate, err := cache.NewPlate(plateType)
	if err != nil {
		t.Fatal(err)
	}

	firstID := wtype.IDOf(firstPlate)

	c := wtype.NewLHComponent(idGen)
	c.Vol = 10.0
	well := firstPlate.GetChildByAddress(wtype.WellCoords{X: 0, Y: 0}).(*wtype.LHWell)
	err = well.SetContents(idGen, c)
	if err != nil {
		t.Fatal(err)
	}

	if firstPlate.IsEmpty(idGen) {
		t.Fatal("Plate shouldn't be empty")
	}

	if err := cache.ReturnPlate(firstPlate); err != nil {
		t.Fatal(err)
	}

	secondPlate, err := cache.NewPlate(plateType)
	if err != nil {
		t.Fatal(err)
	}

	if secondPlate.ID != firstID {
		t.Fatal(fmt.Sprintf("Second plate ID doesn't match the first: %s != %s", secondPlate.ID, firstID))
	}

	if !secondPlate.IsEmpty(idGen) {
		t.Error("secondPlate was not clean")
	}

	thirdPlate, err := cache.NewPlate(plateType)
	if err != nil {
		t.Fatal(err)
	}

	if thirdPlate.ID == firstID {
		t.Fatal("thirdPlate ID was the same as the second, even though the second hasn't been returned")
	}

	if !cache.IsFromCache(secondPlate) {
		t.Error("secondPlate came from cache, but cache.IsFromPlate returned false")
	}

	if !cache.IsFromCache(thirdPlate) {
		t.Error("thirdPlate came from cache, but cache.IsFromPlate returned false")
	}

	fourthPlate, err := inv.PlateTypes.NewPlate(plateType)
	if err != nil {
		t.Fatal(err)
	}

	if cache.IsFromCache(fourthPlate) {
		t.Error("fourthPlate came from inventory, but cache.IsFromPlate returned true")
	}

}
