package plates

import (
	"strings"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

const (
	xStartOffset = 14.28
	yStartOffset = 11.24
	zStartOffset = 0.7
)

func builtinPlates(idGen *id.IDGenerator) composer.Plates {
	plates := makePlates(idGen)
	platesMap := make(composer.Plates)

	for _, p := range plates {
		plateType := composer.PlateType(p.Type)
		cPlate := composer.Plate{
			PlateType:    plateType,
			Manufacturer: p.Mnfr,
			WellShape:    string(p.Welltype.Shape().ShapeName),
			WellH:        p.Welltype.Shape().H,
			WellW:        p.Welltype.Shape().W,
			WellD:        p.Welltype.Shape().D,
			MaxVol:       p.Welltype.MaxVol,
			MinVol:       p.Welltype.Rvol,
			BottomType:   p.Welltype.Bottom,
			BottomH:      p.Welltype.Bottomh,
			WellX:        p.Welltype.Bounds.Size.X,
			WellY:        p.Welltype.Bounds.Size.Y,
			WellZ:        p.Welltype.Bounds.Size.Z,
			ColSize:      p.WellsY(),
			RowSize:      p.WellsX(),
			Height:       p.Height(),
			WellXOffset:  p.WellXOffset,
			WellYOffset:  p.WellYOffset,
			WellXStart:   p.WellXStart,
			WellYStart:   p.WellYStart,
			WellZStart:   p.WellZStart,
			Extra:        p.Welltype.Extra,
		}

		if !strings.Contains(p.Type, "FromSpec") {
			// add offset values to WellX,Y,ZStart
			cPlate = reviseWellStarts(cPlate, xStartOffset, yStartOffset, zStartOffset)
		}
		platesMap[plateType] = cPlate
	}

	return platesMap
}

//		sPlate = reviseWellStarts(sPlate, xStartOffset, yStartOffset, zStartOffset)
func reviseWellStarts(cPlate composer.Plate, xStartOffset, yStartOffset, zStartOffset float64) composer.Plate {
	cPlate.WellXStart += xStartOffset
	cPlate.WellYStart += yStartOffset
	cPlate.WellZStart += zStartOffset

	return cPlate
}
