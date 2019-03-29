package wtype

import (
	"encoding/json"
	"fmt"
)

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
	return tw.dup(false)
}

func (tw *LHTipwaste) DupKeepIDs() *LHTipwaste {
	return tw.dup(true)
}

func (tw *LHTipwaste) dup(keepIDs bool) *LHTipwaste {
	var aw *LHWell
	if keepIDs {
		aw = tw.AsWell.DupKeepIDs()
	} else {
		aw = tw.AsWell.Dup()
	}
	tw2 := NewLHTipwaste(tw.Capacity, tw.Type, tw.Mnfr, tw.Bounds.GetSize(), aw, tw.WellXStart, tw.WellYStart, tw.WellZStart)
	tw2.Contents = tw.Contents
	if keepIDs {
		tw2.ID = tw.ID
		tw2.Name = tw.Name
	}

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

func NewLHTipwaste(capacity int, typ, mfr string, size Coordinates3D, w *LHWell, wellxstart, wellystart, wellzstart float64) *LHTipwaste {
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
	offset := Coordinates3D{
		X: wellxstart - 0.5*w.GetSize().X,
		Y: wellystart - 0.5*w.GetSize().Y,
		Z: wellzstart,
	}
	w.SetOffset(offset) //nolint
	w.Crds = WellCoords{0, 0}

	return &lht
}

func (lht *LHTipwaste) Empty() {
	lht.Contents = 0
}

// Dispose attempt to eject the tips from non-nil channels into the tipwaste.
// Returns a slice of well coordinates which specify where each tip should be ejected (undefined for nil channels),
// and a bool which is true if the tips were disposed of successfully, false if the tipbox is over capacity
func (lht *LHTipwaste) Dispose(channels []*LHChannelParameter) ([]WellCoords, bool) {
	n := 0

	for _, c := range channels {
		if c != nil {
			n += 1
		}
	}

	// currently returning default wellcoords (i.e. A1) for each position is fine, since that's the only position tipboxes have
	// DisponseNum checks that there's space available and increments the contents
	return make([]WellCoords, len(channels)), lht.DisposeNum(n)
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

func (self *LHTipwaste) GetPosition() Coordinates3D {
	if self.parent != nil {
		return self.parent.GetPosition().Add(self.Bounds.GetPosition())
	}
	return self.Bounds.GetPosition()
}

func (self *LHTipwaste) GetSize() Coordinates3D {
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

func (self *LHTipwaste) GetPointIntersections(point Coordinates3D) []LHObject {
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

func (self *LHTipwaste) SetOffset(o Coordinates3D) error {
	self.Bounds.SetPosition(o)
	return nil
}

func (self *LHTipwaste) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *LHTipwaste) ClearParent() {
	self.parent = nil
}

func (self *LHTipwaste) GetParent() LHObject {
	return self.parent
}

//Duplicate copies an LHObject
func (self *LHTipwaste) Duplicate(keepIDs bool) LHObject {
	return self.dup(keepIDs)
}

//DimensionsString returns a string description of the position and size of the object and its children.
func (self *LHTipwaste) DimensionsString() string {
	if self == nil {
		return "nill tipwaste"
	}
	return fmt.Sprintf("Tipwaste \"%s\" at %v+%v\n\t%s", self.GetName(), self.GetPosition(), self.GetSize(), self.AsWell.DimensionsString())
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

func (self *LHTipwaste) CoordsToWellCoords(r Coordinates3D) (WellCoords, Coordinates3D) {
	wc := WellCoords{0, 0}

	c, _ := self.WellCoordsToCoords(wc, TopReference)

	return wc, r.Subtract(c)
}

func (self *LHTipwaste) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates3D, bool) {
	if !self.AddressExists(wc) {
		return Coordinates3D{}, false
	}

	var z float64
	if r == BottomReference {
		z = self.WellZStart
	} else if r == TopReference {
		z = self.WellZStart + self.AsWell.GetSize().Z
	} else {
		return Coordinates3D{}, false
	}

	return self.GetPosition().Add(Coordinates3D{
		self.WellXStart,
		self.WellYStart,
		z}), true
}

//GetTargetOffset get the offset for addressing a well with the named adaptor and channel
func (self *LHTipwaste) GetTargetOffset(adaptorName string, channel int) Coordinates3D {
	targets := self.AsWell.GetWellTargets(adaptorName)
	if channel < 0 || channel >= len(targets) {
		return Coordinates3D{}
	}
	return targets[channel]
}

//GetTargets return all the defined targets for the named adaptor
func (self *LHTipwaste) GetTargets(adaptorName string) []Coordinates3D {
	return self.AsWell.GetWellTargets(adaptorName)
}

func (tw *LHTipwaste) MarshalJSON() ([]byte, error) {
	return json.Marshal(newSTipwaste(tw))
}

func (tw *LHTipwaste) UnmarshalJSON(data []byte) error {
	var stw sTipwaste
	if err := json.Unmarshal(data, &stw); err != nil {
		return err
	}
	stw.Fill(tw)
	return nil
}
