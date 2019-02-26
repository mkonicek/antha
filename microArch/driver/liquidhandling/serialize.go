package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type sProperties struct {
	ID             string
	Positions      map[string]*wtype.LHPosition // position descriptions by position name
	Plates         map[string]*wtype.Plate      // plates by position name
	Tipboxes       map[string]*wtype.LHTipbox   // tipboxes by position name
	Tipwastes      map[string]*wtype.LHTipwaste // tipwastes by position name
	Wastes         map[string]*wtype.Plate      // waste plates by position name
	Washes         map[string]*wtype.Plate      // wash plates by position name
	Model          string
	Mnfr           string
	LHType         LiquidHandlerLevel
	TipType        TipType
	Heads          []*wtype.SerializableHead         // lists every head (whether loaded or not) that is available for the machine
	Adaptors       []*wtype.LHAdaptor                // lists every adaptor (whether loaded or not) that is available for the machine
	HeadAssemblies []*wtype.SerializableHeadAssembly // describes how each loaded head and adaptor is loaded into the machine
	Tips           []*wtype.LHTip
	Preferences    *LayoutOpt
}

func newSProperties(lhp *LHProperties) *sProperties {
	slhp := &sProperties{
		ID:          lhp.ID,
		Positions:   lhp.Positions,
		Plates:      lhp.Plates,
		Tipboxes:    lhp.Tipboxes,
		Tipwastes:   lhp.Tipwastes,
		Wastes:      lhp.Wastes,
		Washes:      lhp.Washes,
		Model:       lhp.Model,
		Mnfr:        lhp.Mnfr,
		LHType:      lhp.LHType,
		TipType:     lhp.TipType,
		Adaptors:    lhp.Adaptors,
		Tips:        lhp.Tips,
		Preferences: lhp.Preferences,
	}

	headIndices := make(map[*wtype.LHHead]int, len(lhp.Heads))
	for i, head := range lhp.Heads {
		headIndices[head] = i
	}
	slhp.HeadAssemblies = make([]*wtype.SerializableHeadAssembly, 0, len(lhp.HeadAssemblies))
	for _, ha := range lhp.HeadAssemblies {
		slhp.HeadAssemblies = append(slhp.HeadAssemblies, wtype.NewSerializableHeadAssembly(ha, headIndices))
	}

	adaptorIndices := make(map[*wtype.LHAdaptor]int, len(lhp.Adaptors))
	for i, adaptor := range lhp.Adaptors {
		adaptorIndices[adaptor] = i
	}
	slhp.Heads = make([]*wtype.SerializableHead, 0, len(lhp.Heads))
	for _, head := range lhp.Heads {
		slhp.Heads = append(slhp.Heads, wtype.NewSerializableHead(head, adaptorIndices))
	}

	return slhp
}

func (slhp *sProperties) Fill(lhp *LHProperties) {
	lhp.ID = slhp.ID
	lhp.Positions = slhp.Positions
	lhp.Plates = slhp.Plates
	lhp.Tipboxes = slhp.Tipboxes
	lhp.Tipwastes = slhp.Tipwastes
	lhp.Wastes = slhp.Wastes
	lhp.Washes = slhp.Washes
	lhp.Model = slhp.Model
	lhp.Mnfr = slhp.Mnfr
	lhp.LHType = slhp.LHType
	lhp.TipType = slhp.TipType
	lhp.Adaptors = slhp.Adaptors
	lhp.Tips = slhp.Tips
	lhp.Preferences = slhp.Preferences

	lhp.Heads = make([]*wtype.LHHead, 0, len(slhp.Heads))
	for _, shead := range slhp.Heads {
		head := wtype.LHHead{}
		shead.Fill(&head, lhp.Adaptors)
		lhp.Heads = append(lhp.Heads, &head)
	}

	lhp.HeadAssemblies = make([]*wtype.LHHeadAssembly, 0, len(slhp.HeadAssemblies))
	for _, sha := range slhp.HeadAssemblies {
		ha := wtype.LHHeadAssembly{}
		sha.Fill(&ha, lhp.Heads)
		lhp.HeadAssemblies = append(lhp.HeadAssemblies, &ha)
	}

	nItems := len(lhp.Plates) + len(lhp.Tipboxes) + len(lhp.Tipwastes) + len(lhp.Wastes) + len(lhp.Washes)
	lhp.PlateLookup = make(map[string]interface{}, nItems)
	lhp.PosLookup = make(map[string]string, nItems)
	lhp.PlateIDLookup = make(map[string]string, nItems)

	for pos, plate := range lhp.Plates {
		lhp.PlateLookup[plate.ID] = plate
		lhp.PosLookup[pos] = plate.ID
		lhp.PlateIDLookup[plate.ID] = pos
	}
	for pos, tipbox := range lhp.Tipboxes {
		lhp.PlateLookup[tipbox.ID] = tipbox
		lhp.PosLookup[pos] = tipbox.ID
		lhp.PlateIDLookup[tipbox.ID] = pos
	}
	for pos, tipwaste := range lhp.Tipwastes {
		lhp.PlateLookup[tipwaste.ID] = tipwaste
		lhp.PosLookup[pos] = tipwaste.ID
		lhp.PlateIDLookup[tipwaste.ID] = pos
	}
	for pos, plate := range lhp.Wastes {
		lhp.PlateLookup[plate.ID] = plate
		lhp.PosLookup[pos] = plate.ID
		lhp.PlateIDLookup[plate.ID] = pos
	}
	for pos, plate := range lhp.Washes {
		lhp.PlateLookup[plate.ID] = plate
		lhp.PosLookup[pos] = plate.ID
		lhp.PlateIDLookup[plate.ID] = pos
	}

}
