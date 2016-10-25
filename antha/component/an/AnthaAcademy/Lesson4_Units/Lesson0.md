### Units:

Units are a critical part of the Antha system. 

The key units you'll need are likely to be:

```
Volume
Mass
Concentration
Temperature
Time
Moles
Rate
Velocity
FlowRate
Length
Area
```

All of these units can be created with a corresponding function from the wunit package following the following format:

```go
wunit.NewVolume(value float64, unitname string)
wunit.NewMass(value float64, unitname string)
wunit.NewConcentration(value float64, unitname string)
```

```go
Vol := wunit.NewVolume(1300, "ul")
Mass := wunit.NewMass(0.001, "g")
Conc := wunit.NewConcentration(0.13, "g/L")
ConcM := wunit.NewConcentration(0.0002, "M/L")
```

All units implement the measurement interface so you can call the following general methods for all units in Antha:


```go
	// the value in base SI units
	SIValue() float64
	// the value in the current units
	RawValue() float64
	// unit plus prefix
	Unit() PrefixedUnit
	// set the value, this must be thread-safe
	// returns old value
	SetValue(v float64) float64
	// convert units
	ConvertTo(p PrefixedUnit) float64
	// wrapper for above
	ConvertToString(s string) float64
	// add to this measurement
	Add(m Measurement)
	// subtract from this measurement
	Subtract(m Measurement)
	// comparison operators
	LessThan(m Measurement) bool
	GreaterThan(m Measurement) bool
	EqualTo(m Measurement) bool
	// A nice string representation
	ToString() string
```

e.g. 

```go
str :=  Vol.ToString()
```

would make str become:

```go
"1300 ul"
```

```go
sivalue :=  Vol.SIValue()
```

would return sivalue as a float64:

```go
0.0013
```

```go
rawvalue :=  Vol.RawValue()
```

would return rawvalue as a float64:

```go
1300
```

