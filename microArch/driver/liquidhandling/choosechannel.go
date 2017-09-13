package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math"
)

// it would probably make more sense for this to be a method of the robot
// in general the instruction generator might well be moved there wholesale
// so that drivers can have specific versions of it... this could lead to some
// very interesting situations though

type ChannelScore float64

type ChannelScoreFunc interface {
	ScoreCombinedChannel(wunit.Volume, *wtype.LHHead, *wtype.LHAdaptor, *wtype.LHTip) ChannelScore
}

type DefaultChannelScoreFunc struct {
}

func (sc DefaultChannelScoreFunc) ScoreCombinedChannel(vol wunit.Volume, head *wtype.LHHead, adaptor *wtype.LHAdaptor, tip *wtype.LHTip) ChannelScore {
	// something pretty simple
	// higher is better
	// 0 == don't bother

	/// NOT USING ADAPTOR??? XXX XXX XXX

	// first we merge the parameters together and see if we can do this at all

	lhcp := head.Params.MergeWithTip(tip)

	// we should always make sure we do not send a volume which is too low

	if lhcp.Minvol.GreaterThanRounded(vol, 1) {
		return 0
	}

	// clearly now vol >= Minvol

	// the main idea is to estimate the error from each source: head, adaptor, tip
	// and make the choice on that basis

	// a big head with a tiny tip will make pretty big errors... a big tip on a tiny
	// head likewise

	// we therefore make our choice as Min(1/tip_error, 1/adaptor_error, 1/head_error)

	err := 999999999.0

	chans := []*wtype.LHChannelParameter{head.GetParams(), tip.GetParams()}

	for _, ch := range chans {
		myerr := sc.ScoreChannel(vol, ch)
		if myerr < err {
			err = myerr
		}
	}

	return ChannelScore(err)
}

func (sc DefaultChannelScoreFunc) ScoreChannel(vol wunit.Volume, lhcp *wtype.LHChannelParameter) float64 {
	// cannot have 0 error
	extra := 1.0
	// we try to estimate the error of using a channel outside its limits
	// first of all how many movements do we need?

	v := vol.RawValue()
	mx := lhcp.Maxvol.ConvertTo(vol.Unit())
	mn := lhcp.Minvol.ConvertTo(vol.Unit())

	n := int(math.Ceil(vol.RawValue() / lhcp.Maxvol.ConvertTo(vol.Unit())))

	// we assume errors scale linearly
	// and that the error is generally greatest at the lowest levels

	tv := v
	if n > 1 {
		tv = mx
	}

	err := (mx-tv)/(mx-mn) + extra

	if n > 1 {
		err *= float64(n + 1)
	}

	score := 1.0 / err

	return score
}

func ChooseChannel(vol wunit.Volume, prms *LHProperties) (*wtype.LHChannelParameter, *wtype.LHTip) {
	var headchosen *wtype.LHHead = nil
	var tipchosen *wtype.LHTip = nil
	var bestscore ChannelScore = ChannelScore(0.0)

	scorer := prms.GetChannelScoreFunc()

	// just choose the best... need to iterate on this sometime though
	// we don't consider head or adaptor changes now

	//fmt.Println("There are ", len(prms.HeadsLoaded), " heads loaded and ", len(prms.Tips), " Tip types available ")

	for _, head := range prms.HeadsLoaded {
		for _, tip := range prms.Tips {
			sc := scorer.ScoreCombinedChannel(vol, head, head.Adaptor, tip)
			if sc > bestscore {
				headchosen = head
				tipchosen = tip
				bestscore = sc
			}
		}

	}

	if headchosen == nil {
		return nil, nil
	}

	// shouldn't we also return the adaptor?
	// and probably the whole head rather than just its channel parameters

	return headchosen.GetParams(), tipchosen
}

func ChooseChannels(vols []wunit.Volume, prms *LHProperties) ([]*wtype.LHChannelParameter, []*wtype.LHTip, []string, error) {
	prmA := make([]*wtype.LHChannelParameter, len(vols))
	tipA := make([]*wtype.LHTip, len(vols))
	tipTypeA := make([]string, len(vols))

	// we choose individually

	for i := 0; i < len(vols); i++ {
		if vols[i].IsZero() {
			continue
		}
		prm, tip := ChooseChannel(vols[i], prms)
		if tip == nil {
			return prmA, tipA, tipTypeA, fmt.Errorf(TipChosenError(vols[i], prms))
		}
		prmA[i] = prm
		tipA[i] = tip.Dup()
		tipTypeA[i] = tip.Type
	}

	return prmA, tipA, tipTypeA, nil
}
