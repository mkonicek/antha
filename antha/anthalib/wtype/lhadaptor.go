package wtype

// adaptor
type LHAdaptor struct {
	Name         string
	ID           string
	Manufacturer string
	Params       *LHChannelParameter
	Tips         []*LHTip
}

//NewLHAdaptor make a new adaptor
func NewLHAdaptor(name, mf string, params *LHChannelParameter) *LHAdaptor {
	return &LHAdaptor{
		Name:         name,
		Manufacturer: mf,
		Params:       params,
		ID:           GetUUID(),
		Tips:         make([]*LHTip, params.Multi),
	}
}

//Dup duplicate the adaptor and any loaded tips with new IDs
func (lha *LHAdaptor) Dup() *LHAdaptor {
	return lha.dup(false)
}

//AdaptorType the manufacturer and name of the adaptor
func (lha *LHAdaptor) AdaptorType() string {
	return lha.Manufacturer + lha.Name
}

//DupKeepIDs duplicate the adaptor and any loaded tips keeping the same IDs
func (lha *LHAdaptor) DupKeepIDs() *LHAdaptor {
	return lha.dup(true)
}

func (lha *LHAdaptor) dup(keepIDs bool) *LHAdaptor {

	var ad *LHAdaptor
	if keepIDs {
		ad = NewLHAdaptor(lha.Name, lha.Manufacturer, lha.Params.DupKeepIDs())
		ad.ID = lha.ID
		for i, tip := range lha.Tips {
			ad.AddTip(i, tip.DupKeepID())
		}
	} else {
		ad = NewLHAdaptor(lha.Name, lha.Manufacturer, lha.Params.Dup())
		for i, tip := range lha.Tips {
			ad.AddTip(i, tip.Dup())
		}
	}

	return ad
}

//NumTipsLoaded the number of tips currently loaded
func (lha *LHAdaptor) NumTipsLoaded() int {
	r := 0
	for _, tip := range lha.Tips {
		if tip != nil {
			r += 1
		}
	}
	return r
}

//IsTipLoaded Is there a tip loaded on channelNumber
func (lha *LHAdaptor) IsTipLoaded(channelNumber int) bool {
	return lha.Tips[channelNumber] != nil
}

//GetTip Return the tip at channelNumber, nil otherwise
func (lha *LHAdaptor) GetTip(channelNumber int) *LHTip {
	return lha.Tips[channelNumber]
}

//AddTip Load a tip to the specified channel, overwriting any tip already present
func (lha *LHAdaptor) AddTip(channelNumber int, tip *LHTip) {
	lha.Tips[channelNumber] = tip
}

//RemoveTip Remove a tip from the specified channel and return it
func (lha *LHAdaptor) RemoveTip(channelNumber int) *LHTip {
	tip := lha.Tips[channelNumber]
	lha.Tips[channelNumber] = nil
	return tip
}

//RemoveTips Return all previously loaded tips, with nils removed
func (lha *LHAdaptor) RemoveTips() []*LHTip {
	ret := make([]*LHTip, 0, lha.NumTipsLoaded())
	for _, tip := range lha.Tips {
		if tip != nil {
			ret = append(ret, tip)
		}
	}
	lha.Tips = make([]*LHTip, lha.Params.Multi)
	return ret
}

//GetParams get the channel parameters for the adaptor, combined with any loaded tips
func (lha *LHAdaptor) GetParams() *LHChannelParameter {
	if lha.NumTipsLoaded() == 0 {
		return lha.Params
	} else {
		params := *lha.Params
		for _, tip := range lha.Tips {
			if tip != nil {
				params = *params.MergeWithTip(tip)
			}
		}
		return &params
	}
}

//GetSmallestChannelSpacing get the smallest possible distance between successive channels
func (self *LHAdaptor) GetSmallestChannelSpacing() Coordinates2D {

	//hjk: currently assume a constant fixed offset between channels
	//     this will need updating when we support better reporting of head
	//     capabilities wrt. independent multi-channelling
	channelStep := 9.0 //9mm between channels
	if self.Params.Orientation == LHVChannel {
		return Coordinates2D{Y: channelStep}
	}
	return Coordinates2D{X: channelStep}
}

//GetLargestChannelSpacing get the largest possible distance between successive channels
func (self *LHAdaptor) GetLargestChannelSpacing() Coordinates2D {
	//equal to smallest if independent
	if !self.Params.Independent {
		return self.GetSmallestChannelSpacing()
	}

	//completely arbitrary for now since we don't report this
	return self.GetSmallestChannelSpacing().Multiply(2.0)
}

//GetMostCompactChannelPositions get the relative channel positions for the adaptor
//in the most tightly bunched layout supported
func (self *LHAdaptor) GetMostCompactChannelPositions() ChannelPositions {
	ret := make([]Coordinates2D, self.Params.Multi)
	offset := self.GetSmallestChannelSpacing()
	current := Coordinates2D{}

	for i := range ret {
		ret[i] = current
		current = current.Add(offset)
	}

	return ret
}

//GetWellTargets get the offset from the center of the well for each channel that
//can access the well simultaneously
//returns nil if the well cannot fit multiple channels
func (self *LHAdaptor) GetWellTargets(well *LHWell) []Coordinates2D {
	//no well targets if there's no adaptor loaded
	if self == nil {
		return nil
	}
	channelPositions := self.GetMostCompactChannelPositions()
	channelRadius := 3.0 //should come from adaptor

	//total size of channels, including radius
	channelSize := channelPositions.Size(channelRadius)

	if wellSize := well.GetSize(); wellSize.X < channelSize.X || wellSize.Y < channelSize.Y {
		return nil
	}

	//set the channel positions center as their origin
	return channelPositions.Subtract(channelPositions.GetCenter())
}

//A list of 2d coordinates of the channels of an adaptor in channel order
type ChannelPositions []Coordinates2D

//Size get the total footprint size of the channel positions including the radius
func (self ChannelPositions) Size(channelRadius float64) Coordinates2D {
	if len(self) == 0 {
		return Coordinates2D{}
	}
	return Coordinates2D{
		X: self[len(self)-1].X - self[0].X + 2*channelRadius,
		Y: self[len(self)-1].Y - self[0].Y + 2*channelRadius,
	}
}

//Add return a new set of channel positions with the offset added
func (self ChannelPositions) Add(rhs Coordinates2D) ChannelPositions {
	ret := make(ChannelPositions, len(self))
	for i, crd := range self {
		ret[i] = crd.Add(rhs)
	}
	return ret
}

//Subtract return a new set of channel positions with the offset subtracted
func (self ChannelPositions) Subtract(rhs Coordinates2D) ChannelPositions {
	ret := make(ChannelPositions, len(self))
	for i, crd := range self {
		ret[i] = crd.Subtract(rhs)
	}
	return ret
}

//GetCenter return the center of the channel coordinates
func (self ChannelPositions) GetCenter() Coordinates2D {
	return NewBoundingRectangle(self).Center()
}
