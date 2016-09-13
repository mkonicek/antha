package ast

import "fmt"

// TODO make more specific
type Location interface{}

type Movement struct {
	From Location
	To   Location
}

type Request struct {
	MixVol *Interval
	Temp   *Interval
	Time   *Interval
	Move   []Movement
	Manual bool
}

func (a Request) String() string {
	var r []string
	if a.MixVol != nil {
		a, b := a.MixVol.Extrema()
		r = append(r, fmt.Sprintf("mix %fl %fl", a, b))
	}
	if a.Temp != nil {
		a, b := a.Temp.Extrema()
		r = append(r, fmt.Sprintf("hold %fC %fC", a, b))
	}
	if a.Time != nil {
		a, b := a.Time.Extrema()
		r = append(r, fmt.Sprintf("hold %fs %fs", a, b))
	}
	for _, m := range a.Move {
		r = append(r, fmt.Sprintf("move %s %s", m.From, m.To))
	}
	if a.Manual {
		r = append(r, "manual")
	}
	return fmt.Sprint(r)
}

// Compute greatest lower bound of a set of requests
func Meet(reqs []Request) (req Request) {
	meetI := func(dst **Interval, src *Interval) {
		if *dst == nil {
			*dst = src
		} else if src != nil {
			(*dst).Add(*src)
		}
	}
	for _, r := range reqs {
		meetI(&req.MixVol, r.MixVol)
		meetI(&req.Temp, r.Temp)
		meetI(&req.Time, r.Time)
		req.Move = append(req.Move, r.Move...)
		if !req.Manual && r.Manual {
			req.Manual = r.Manual
		}
	}
	return
}
