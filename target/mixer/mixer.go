package mixer

/*
func (a *Mixer) saveFile(name string) ([]byte, error) {
	data, status := a.driver.GetOutputFile()
	if err := status.GetError(); err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	bs := []byte(data)

	if err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    0644,
		Size:    int64(len(bs)),
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	} else if _, err := tw.Write(bs); err != nil {
		return nil, err
	} else if err := tw.Close(); err != nil {
		return nil, err
	} else if err := gw.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (a *Mixer) makeMix(ctx context.Context, mixes []*wtype.LHInstruction) (*target.Mix, error) {

	err = r.Liquidhandler.MakeSolutions(ctx, r.LHRequest)
	// TODO: MIS unfortunately we need to make sure this stays up to date would
	// be better to remove this and just use the ones the liquid handler holds
	r.LHProperties = r.Liquidhandler.Properties

	if err != nil {
		return nil, err
	}

	name := a.opt.DriverOutputFileName
	if len(name) == 0 {
		// TODO: Desired filename not exposed in current driver interface, so pick
		// a name. So far, at least Gilson software cares what the filename is, so
		// use .sqlite for compatibility
		name = strings.Replace(fmt.Sprintf("%s.sqlite", time.Now().Format(time.RFC3339)), ":", "_", -1)
	}

	tarball, err := a.saveFile(name)
	if err != nil {
		return nil, err
	}

	return &target.Mix{
		Dev:             a,
		Request:         r.LHRequest,
		Properties:      r.LHProperties,
		FinalProperties: r.Liquidhandler.FinalProperties,
		Final:           r.Liquidhandler.PlateIDMap(),
		Files: target.Files{
			Tarball: tarball,
			Type:    a.FileType(),
		},
	}, nil
}

// New creates a new Mixer
func New(opt Opt, d driver.LiquidhandlingDriver) (*Mixer, error) {
	userPreferences := &driver.LayoutOpt{
		Tipboxes:  driver.Addresses(opt.DriverSpecificTipPreferences),
		Inputs:    driver.Addresses(opt.DriverSpecificInputPreferences),
		Outputs:   driver.Addresses(opt.DriverSpecificOutputPreferences),
		Tipwastes: driver.Addresses(opt.DriverSpecificTipWastePreferences),
		Washes:    driver.Addresses(opt.DriverSpecificWashPreferences),
	}

	if p, status := d.GetCapabilities(); !status.Ok() {
		return nil, status.GetError()
	} else if err := p.ApplyUserPreferences(userPreferences); err != nil {
		return nil, err
	} else {
		p.Driver = d
		return &Mixer{driver: d, properties: &p, opt: opt}, nil
	}
}
*/
