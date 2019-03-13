package plateCache

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

const plateType = "pcrplate"

func makeContext() context.Context {
	ctx := NewContext(testinventory.NewContext(context.Background()))
	return ctx
}

func TestPlateReuse(t *testing.T) {

	ctx := makeContext()
	firstPlate, err := cache.NewPlate(ctx, plateType)
	if err != nil {
		t.Fatal(err)
	}

	firstID := wtype.IDOf(firstPlate)

	c := wtype.NewLHComponent()
	c.Vol = 10.0
	well := firstPlate.GetChildByAddress(wtype.WellCoords{X: 0, Y: 0}).(*wtype.LHWell)
	err = well.SetContents(c)
	if err != nil {
		t.Fatal(err)
	}

	if firstPlate.IsEmpty() {
		t.Fatal("Plate shouldn't be empty")
	}

	if err := cache.ReturnObject(ctx, firstPlate); err != nil {
		t.Fatal(err)
	}

	secondPlate, err := cache.NewPlate(ctx, plateType)
	if err != nil {
		t.Fatal(err)
	}

	if secondPlate.ID != firstID {
		t.Fatal(fmt.Sprintf("Second plate ID doesn't match the first: %s != %s", secondPlate.ID, firstID))
	}

	if !secondPlate.IsEmpty() {
		t.Error("secondPlate was not clean")
	}

	thirdPlate, err := cache.NewPlate(ctx, plateType)
	if err != nil {
		t.Fatal(err)
	}

	if thirdPlate.ID == firstID {
		t.Fatal("thirdPlate ID was the same as the second, even though the second hasn't been returned")
	}

	if !cache.IsFromCache(ctx, secondPlate) {
		t.Error("secondPlate came from cache, but cache.IsFromPlate returned false")
	}

	if !cache.IsFromCache(ctx, thirdPlate) {
		t.Error("thirdPlate came from cache, but cache.IsFromPlate returned false")
	}

	fourthPlate, err := inventory.NewPlate(ctx, plateType)
	if err != nil {
		t.Fatal(err)
	}

	if cache.IsFromCache(ctx, fourthPlate) {
		t.Error("fourthPlate came from inventory, but cache.IsFromPlate returned true")
	}

}
