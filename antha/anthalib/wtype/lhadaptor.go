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
		Tips:         []*LHTip{},
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
	var params *LHChannelParameter
	if keepIDs {
		params = lha.Params.DupKeepIDs()
	} else {
		params = lha.Params.Dup()
	}

	ad := NewLHAdaptor(lha.Name, lha.Manufacturer, params)

	for i, tip := range lha.Tips {
		if tip != nil {
			if keepIDs {
				ad.AddTip(i, tip.DupKeepID())
			} else {
				ad.AddTip(i, tip.Dup())
			}
		}
	}

	if keepIDs {
		ad.ID = lha.ID
	} else {
		ad.ID = GetUUID()
	}

	return ad
}

//NTipsLoaded the number of tips currently loaded
func (lha *LHAdaptor) NTipsLoaded() int {
	r := 0
	for i := range lha.Tips {
		if lha.Tips[i] != nil {
			r += 1
		}
	}
	return r
}

//IsTipLoaded Is there a tip loaded on channel_number
func (lha *LHAdaptor) IsTipLoaded(channel_number int) bool {
	return lha.Tips[channel_number] != nil
}

//GetTip Return the tip at channel_number, nil otherwise
func (lha *LHAdaptor) GetTip(channel_number int) *LHTip {
	return lha.Tips[channel_number]
}

//AddTip Load a tip to the specified channel
func (lha *LHAdaptor) AddTip(channel_number int, tip *LHTip) {
	lha.Tips[channel_number] = tip
}

//RemoveTip Remove a tip from the specified channel and return it
func (lha *LHAdaptor) RemoveTip(channel_number int) *LHTip {
	tip := lha.Tips[channel_number]
	lha.Tips[channel_number] = nil
	return tip
}

//RemoveTips Remove every tip from the adaptor
func (lha *LHAdaptor) RemoveTips() []*LHTip {
	ret := make([]*LHTip, 0, lha.NTipsLoaded())
	for i := range lha.Tips {
		if lha.Tips[i] != nil {
			ret = append(ret, lha.Tips[i])
			lha.Tips[i] = nil
		}
	}
	return ret
}

//GetParams get the channel parameters for the adaptor, combined with any loaded tips
func (lha *LHAdaptor) GetParams() *LHChannelParameter {
	if lha.NTipsLoaded() == 0 {
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

//GetSmallestChannelOffset get the smallest possible distance between successive channels
func (self *LHAdaptor) GetSmallestChannelOffset() Coordinates2D {

	//hjk: currently assume a constant fixed offset between channels
	//     this will need updating when we support better reporting of head
	//     capabilities wrt. independent multi-channelling
	channelStep := 9.0 //9mm between channels
	if self.Params.Orientation == LHVChannel {
		return Coordinates2D{Y: channelStep}
	}
	return Coordinates2D{X: channelStep}
}

//GetLargestChannelOffset get the largest possible distance between successive channels
func (self *LHAdaptor) GetLargestChannelOffset() Coordinates2D {
	//equal to smallest if independent
	if !self.Params.Independent {
		return self.GetSmallestChannelOffset()
	}

	//completely arbitrary for now since we don't report this
	return self.GetSmallestChannelOffset().Multiply(2.0)
}

//GetMostCompactChannelPositions get the relative channel positions for the adaptor
//in the most tightly bunched layout supported
func (self *LHAdaptor) GetMostCompactChannelPositions() ChannelPositions {
	ret := make([]Coordinates2D, self.Params.Multi)
	offset := self.GetSmallestChannelOffset()
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
	channelPositions = channelPositions.Subtract(channelPositions.GetCenter())

	return channelPositions
}

//A list of 2d coordinates of the channels of an adaptor in channel order
type ChannelPositions []Coordinates2D

//Size get the total footprint size of the channel positions including the radius
func (self ChannelPositions) Size(channelRadius float64) Coordinates2D {
	if len(self) <= 1 {
		return Coordinates2D{}
	}
	return Coordinates2D{
		X: self[len(self)-1].X - self[0].X + 2*channelRadius,
		Y: self[len(self)-1].Y - self[0].Y + 2*channelRadius,
	}
}

//Add return a new set of channel positions with the offset added
func (self ChannelPositions) Add(rhs Coordinates2D) ChannelPositions {
	ret := make([]Coordinates2D, len(self))
	for i, crd := range self {
		ret[i] = crd.Add(rhs)
	}
	return ret
}

//Subtract return a new set of channel positions with the offset subtracted
func (self ChannelPositions) Subtract(rhs Coordinates2D) ChannelPositions {
	ret := make([]Coordinates2D, len(self))
	for i, crd := range self {
		ret[i] = crd.Subtract(rhs)
	}
	return ret
}

//GetCenter return the mean of the channel coordinates
func (self ChannelPositions) GetCenter() Coordinates2D {
	ret := Coordinates2D{}
	for _, crd := range self {
		ret = ret.Add(crd)
	}
	return ret.Divide(float64(len(self)))
}
