package Optimization

type AssemblyMethod interface {
	Name() string                                             // what is this method called
	Type() string                                             // method classification
	CanAssemble([]string) bool                                // can these parts work?
	Validate([]string, string) (bool, error)                  // do these parts make this output?
	Assemble([]string) ([]string, error)                      // what do these parts make?
	Disassemble(string, AssemblyParameters) ([]string, error) // split this sequence into parts
	RawToParts([]string) ([]string, error)                    // convert each part
	PartsToRaw([]string) ([]string, error)                    // return each part to its raw state
}

// e.g.			      Assemble()          RawToParts()
//                         <--------------       <-------------
// sequence + split points                 parts                raw fragments
//                         -------------->       ------------->
//                          Disassemble()         PartsToRaw()

type AssemblyParameters interface {
	SplitPoints() []int
}
