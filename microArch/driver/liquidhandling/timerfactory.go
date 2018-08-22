package liquidhandling

import (
	"fmt"
	"time"
)

// urrgh -- this needs to get packaged in with the driver

func GetTimerFor(model, mnfr string) LHTimer {
	timers := makeTimers()

	fmt.Println("Getting timer for ", model+mnfr)

	_, ok := timers[model+mnfr]
	if ok {
		return timers[model+mnfr]
	} else {
		fmt.Println("None found")
		return makeNullTimer()
	}
}

func makeTimers() map[string]LHTimer {
	timers := make(map[string]LHTimer, 2)
	timers["GilsonPipetmax"] = makeGilsonPipetmaxTimer()
	timers["CyBioFelix"] = makeCyBioFelixTimer()
	timers["CyBioGeneTheatre"] = makeCyBioGeneTheatreTimer()
	timers["LabcyteEcho550"] = makeLabcyteEchoTimer("550")
	timers["LabcyteEcho520"] = makeLabcyteEchoTimer("520")
	return timers
}

func makeNullTimer() LHTimer {
	// always returns zero
	t := NewTimer()

	return t
}

func makeGilsonPipetmaxTimer() LHTimer {
	t := NewTimer()
	t.Times[INI], _ = time.ParseDuration("5s") // INI
	t.Times[LDT], _ = time.ParseDuration("7s") // LDT
	t.Times[UDT], _ = time.ParseDuration("7s") // UDT
	t.Times[SUK], _ = time.ParseDuration("4s") // SUK
	t.Times[BLW], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[ASP], _ = time.ParseDuration("0.8s") // ASP
	t.Times[DSP], _ = time.ParseDuration("0.8s") // DSP
	t.Times[BLO], _ = time.ParseDuration("0.7s") // BLO
	t.Times[PTZ], _ = time.ParseDuration("0s")   // PTZ
	t.Times[MOV], _ = time.ParseDuration("2.3s") // MOV	-- using mean figures for horizontal (0.94) and 2x vertical (1.36)
	t.Times[LOD], _ = time.ParseDuration("5.6s") // LOAD
	t.Times[ULD], _ = time.ParseDuration("5.4s") // UNLOAD
	t.Times[MIX], _ = time.ParseDuration("1.5s") // MIX

	return t
}

func makeCyBioFelixTimer() LHTimer {
	t := NewTimer()
	t.Times[LDT], _ = time.ParseDuration("8s") // LDT
	t.Times[UDT], _ = time.ParseDuration("6s") // UDT
	t.Times[SUK], _ = time.ParseDuration("4s") // SUK
	t.Times[BLW], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[ASP], _ = time.ParseDuration("12s")  // ASP
	t.Times[DSP], _ = time.ParseDuration("10s")  // DSP
	t.Times[BLO], _ = time.ParseDuration("10s")  // BLO
	t.Times[PTZ], _ = time.ParseDuration("0.5s") // PTZ
	t.Times[MOV], _ = time.ParseDuration("0s")   // MOV
	t.Times[LOD], _ = time.ParseDuration("10s")  // LOAD
	t.Times[ULD], _ = time.ParseDuration("12s")  // UNLOAD
	t.Times[MIX], _ = time.ParseDuration("28s")  // MIX

	return t
}

func makeCyBioGeneTheatreTimer() LHTimer {
	t := NewTimer()
	t.Times[LDT], _ = time.ParseDuration("8s") // LDT
	t.Times[UDT], _ = time.ParseDuration("6s") // UDT
	t.Times[SUK], _ = time.ParseDuration("4s") // SUK
	t.Times[BLW], _ = time.ParseDuration("4s") // BLW

	// lower level instructions

	t.Times[ASP], _ = time.ParseDuration("9s")   // ASP
	t.Times[DSP], _ = time.ParseDuration("10s")  // DSP
	t.Times[BLO], _ = time.ParseDuration("10s")  // BLO
	t.Times[PTZ], _ = time.ParseDuration("0.5s") // PTZ
	t.Times[MOV], _ = time.ParseDuration("0s")   // MOV
	t.Times[LOD], _ = time.ParseDuration("10s")  // LOAD
	t.Times[ULD], _ = time.ParseDuration("12s")  // UNLOAD
	t.Times[MIX], _ = time.ParseDuration("13s")  // MIX

	return t
}

func makeLabcyteEchoTimer(model string) LHTimer {
	flowRate := 5000.0 // nL / s make model dependent
	moveRate := 0.1    // s/well ditto... also find out what this actually is
	scanRate := 0.1    // s/well  ditto

	timer := highLeveltimer{name: fmt.Sprintf("LabcyteEcho%s", model), model: model, flowRate: flowRate, moveRate: moveRate, scanRate: scanRate}

	return timer
}
