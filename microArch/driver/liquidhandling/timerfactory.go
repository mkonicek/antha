package liquidhandling

import "time"

// urrgh -- this needs to get packaged in with the driver

func GetTimerFor(model, mnfr string) *LHTimer {
	timers := makeTimers()

	//fmt.Println("Getting timer for ", model, mnfr)

	_, ok := timers[model+mnfr]
	if ok {
		return timers[model+mnfr]
	} else {
		//fmt.Println("None found")
		return makeNullTimer()
	}

	return nil
}

func makeTimers() map[string]*LHTimer {
	timers := make(map[string]*LHTimer, 2)

	timers["GilsonPipetmax"] = makeGilsonPipetmaxTimer()
	timers["CyBioFelix"] = makeCyBioFelixTimer()
	timers["CyBioGeneTheatre"] = makeCyBioGeneTheatreTimer()
	return timers
}

func makeNullTimer() *LHTimer {
	// always returns zero
	t := NewTimer()

	return t
}

func makeGilsonPipetmaxTimer() *LHTimer {
	t := NewTimer()
	t.Times[7], _ = time.ParseDuration("8s")  // LDT
	t.Times[8], _ = time.ParseDuration("6s")  // UDT
	t.Times[19], _ = time.ParseDuration("4s") // SUK
	t.Times[20], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[11], _ = time.ParseDuration("4s")   // ASP
	t.Times[12], _ = time.ParseDuration("4s")   // DSP
	t.Times[13], _ = time.ParseDuration("4s")   // BLO
	t.Times[14], _ = time.ParseDuration("0.5s") // PTZ
	t.Times[15], _ = time.ParseDuration("2s")   // MOV
	t.Times[17], _ = time.ParseDuration("6s")   // LOAD
	t.Times[18], _ = time.ParseDuration("8s")   // UNLOAD
	t.Times[32], _ = time.ParseDuration("6s")   // MIX

	return t
}

func makeCyBioFelixTimer() *LHTimer {
	t := NewTimer()
	t.Times[7], _ = time.ParseDuration("8s")  // LDT
	t.Times[8], _ = time.ParseDuration("6s")  // UDT
	t.Times[19], _ = time.ParseDuration("4s") // SUK
	t.Times[20], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[11], _ = time.ParseDuration("12s")  // ASP
	t.Times[12], _ = time.ParseDuration("10s")  // DSP
	t.Times[13], _ = time.ParseDuration("10s")  // BLO
	t.Times[14], _ = time.ParseDuration("0.5s") // PTZ
	t.Times[15], _ = time.ParseDuration("0s")   // MOV
	t.Times[17], _ = time.ParseDuration("10s")  // LOAD
	t.Times[18], _ = time.ParseDuration("12s")  // UNLOAD
	t.Times[32], _ = time.ParseDuration("28s")  // MIX

	return t
}

func makeCyBioGeneTheatreTimer() *LHTimer {
	t := NewTimer()
	t.Times[7], _ = time.ParseDuration("8s")  // LDT
	t.Times[8], _ = time.ParseDuration("6s")  // UDT
	t.Times[19], _ = time.ParseDuration("4s") // SUK
	t.Times[20], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[11], _ = time.ParseDuration("9s")   // ASP
	t.Times[12], _ = time.ParseDuration("10s")  // DSP
	t.Times[13], _ = time.ParseDuration("10s")  // BLO
	t.Times[14], _ = time.ParseDuration("0.5s") // PTZ
	t.Times[15], _ = time.ParseDuration("0s")   // MOV
	t.Times[17], _ = time.ParseDuration("10s")  // LOAD
	t.Times[18], _ = time.ParseDuration("12s")  // UNLOAD
	t.Times[32], _ = time.ParseDuration("13s")  // MIX

	return t
}
