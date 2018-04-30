package wtype

// defines a tip waste
import "fmt"

// tip waste

type LHTipwaste struct {
	Name       string
	ID         string
	Type       string
	Mnfr       string
	Capacity   int
	Contents   int
	Height     float64
	WellXStart float64
	WellYStart float64
	WellZStart float64
	AsWell     *LHWell
	Bounds     BBox
	parent     LHObject `gotopb:"-"`
}

func (tw LHTipwaste) SpaceLeft() int {
	return tw.Capacity - tw.Contents
}

func (te LHTipwaste) String() string {
	return fmt.Sprintf(
		`LHTipwaste {
	ID: %s,
	Type: %s,
    Name: %s,
	Mnfr: %s,
	Capacity: %d,
	Contents: %d,
    Length: %f,
    Width: %f,
	Height: %f,
	WellXStart: %f,
	WellYStart: %f,
	WellZStart: %f,
	AsWell: %p,
}
`,
		te.ID,
		te.Type,
		te.Name,
		te.Mnfr,
		te.Capacity,
		te.Contents,
		te.Bounds.GetSize().X,
		te.Bounds.GetSize().Y,
		te.Bounds.GetSize().Z,
		te.WellXStart,
		te.WellYStart,
		te.WellZStart,
		te.AsWell, //AsWell is printed as pointer to keep things short
	)
}

func (tw *LHTipwaste) Dup() *LHTipwaste {
	tw2 := NewLHTipwaste(tw.Capacity, tw.Type, tw.Mnfr, tw.Bounds.GetSize(), tw.AsWell, tw.WellXStart, tw.WellYStart, tw.WellZStart)

	tw2.Contents = tw.Contents

	return tw2
}

func (tw *LHTipwaste) GetName() string {
	if tw == nil {
		return "<nil>"
	}
	return tw.Name
}

func (tw *LHTipwaste) GetID() string {
	return tw.ID
}

func (tw *LHTipwaste) GetType() string {
	if tw == nil {
		return "<nil>"
	}
	return tw.Type
}

func (self *LHTipwaste) GetClass() string {
	return "tipwaste"
}

func NewLHTipwaste(capacity int, typ, mfr string, size Coordinates, w *LHWell, wellxstart, wellystart, wellzstart float64) *LHTipwaste {
	var lht LHTipwaste
	//	lht.ID = "tipwaste-" + GetUUID()
	lht.ID = GetUUID()
	lht.Type = typ
	lht.Name = fmt.Sprintf("%s_%s", typ, lht.ID[1:len(lht.ID)-2])
	lht.Mnfr = mfr
	lht.Capacity = capacity
	lht.Bounds.SetSize(size)
	lht.AsWell = w
	lht.WellXStart = wellxstart
	lht.WellYStart = wellystart
	lht.WellZStart = wellzstart

	w.SetParent(&lht) //nolint

	return &lht
}

func (lht *LHTipwaste) Empty() {
	lht.Contents = 0
}

func (lht *LHTipwaste) Dispose(channels []*LHChannelParameter) bool {
	// this just checks numbers for now
	n := 0

	for _, c := range channels {
		if c != nil {
			n += 1
		}
	}

	return lht.DisposeNum(n)
}

func (lht *LHTipwaste) DisposeNum(num int) bool {
	if lht.Capacity-lht.Contents < num {
		return false
	}

	lht.Contents += num
	return true

}

//##############################################
//@implement LHObject
//##############################################

func (self *LHTipwaste) GetPosition() Coordinates {
	if self.parent != nil {
		return self.parent.GetPosition().Add(self.Bounds.GetPosition())
	}
	return self.Bounds.GetPosition()
}

func (self *LHTipwaste) GetSize() Coordinates {
	return self.Bounds.GetSize()
}

func (self *LHTipwaste) GetBoxIntersections(box BBox) []LHObject {
	if r := self.AsWell.GetBoxIntersections(box); len(r) > 0 {
		return r
	}

	ret := []LHObject{}
	//relative box
	box.SetPosition(box.GetPosition().Subtract(OriginOf(self)))
	if self.Bounds.IntersectsBox(box) {
		ret = append(ret, self)
	}
	return ret
}

func (self *LHTipwaste) GetPointIntersections(point Coordinates) []LHObject {
	if r := self.AsWell.GetPointIntersections(point); len(r) > 0 {
		return r
	}

	//relative point
	point = point.Subtract(OriginOf(self))

	ret := []LHObject{}
	//Todo, test well
	if self.Bounds.IntersectsPoint(point) {
		ret = append(ret, self)
	}
	return ret
}

func (self *LHTipwaste) SetOffset(o Coordinates) error {
	self.Bounds.SetPosition(o)
	return nil
}

func (self *LHTipwaste) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

func (self *LHTipwaste) GetParent() LHObject {
	return self.parent
}

//##############################################
//@implement Addressable
//##############################################

func (self *LHTipwaste) AddressExists(c WellCoords) bool {
	return c.X == 0 && c.Y == 0
}

func (self *LHTipwaste) NRows() int {
	return 1
}

func (self *LHTipwaste) NCols() int {
	return 1
}

func (self *LHTipwaste) GetChildByAddress(c WellCoords) LHObject {
	if !self.AddressExists(c) {
		return nil
	}
	//LHWells arent LHObjects yet
	return self.AsWell
}

func (self *LHTipwaste) CoordsToWellCoords(r Coordinates) (WellCoords, Coordinates) {
	wc := WellCoords{0, 0}

	c, _ := self.WellCoordsToCoords(wc, TopReference)

	return wc, r.Subtract(c)
}

func (self *LHTipwaste) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
	if !self.AddressExists(wc) {
		return Coordinates{}, false
	}

	var z float64
	if r == BottomReference {
		z = self.WellZStart
	} else if r == TopReference {
		z = self.WellZStart + self.AsWell.GetSize().Z
	} else {
		return Coordinates{}, false
	}

	return self.GetPosition().Add(Coordinates{
		self.WellXStart + 0.5*self.AsWell.GetSize().X,
		self.WellYStart + 0.5*self.AsWell.GetSize().Y,
		z}), true
}
