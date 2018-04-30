package liquidhandling

// no longer need to supply tipboxes after the fact
func (lh *Liquidhandler) Refresh_tipboxes_tipwastes(rq *LHRequest) {

	// dead simple

	lh.FinalProperties.RemoveTipBoxes()

	for pos := range lh.Properties.PosLookup {
		tb, ok := lh.Properties.Tipboxes[pos]

		if ok {
			lh.FinalProperties.AddTipBoxTo(pos, tb.Dup())
			tb.Refresh()
			continue
		}

		tw, ok := lh.Properties.Tipwastes[pos]

		if ok {
			// swap the wastes
			tw2 := lh.FinalProperties.Tipwastes[pos]
			tw2.Contents = tw.Contents
			tw.Empty()
		}
	}
}
