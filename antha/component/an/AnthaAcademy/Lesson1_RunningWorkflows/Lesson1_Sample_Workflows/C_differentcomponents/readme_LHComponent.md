### LHComponents:

One of the key antha types which will typically be specified in the parameters file is the wtype.LHComponent

LHComponents can be accessed in the parameters.yml file in the following way:

These are written as a string: e.g. 
"Diluent":"water",
“dnastock”:”gfpstock”,

Before a component can be used, currently, the concept of that component needs to be added to the factory.
i.e. When we say the concept of a component we don't mean a specific sample of water, which would be called from an inventory instead, but any sample of water, i.e. which has the liquidhandling properties of water.

#### Checking available components from command line
To check what components are available type the following command

```bash
antharun lhcomponents
```

## Factory
The factory is located in the following path:

```bash
$GOPATH/src/github.com/antha-lang/antha/microArch/factory/
```

Here you can find both the plate factory and component factory

### Component factory:

Open the file and add the component to the list within the body of the func makeComponentLibrary()

e.g.

```go
A = wtype.NewLHComponent()
	A.CName = "tartrazine"
	A.Type = wtype.LTWater // or could use wtype.LiquidTypeFromString("water")
	A.Smax = 9999
	cmap[A.CName] = A
```

therefore a new component would be specified as follows:

```go
A = wtype.NewLHComponent()
    A.CName = "mynewviscouscomponent"
    A.Type = wtype.LTVISCOUS
    A.Smax = 9999
    cmap[A.CName] = A
```

### plate factory:

Open the file and add the component to the list within the body of the func makePlateLibrary()


#### Checking available plates from command line
To check what components are available type the following command

```bash
antharun lhplates
```

### LiquidTypes:
	
You may want to change the .Type to something else as this will determine how the liquid type is pipetted. 
Currently this consists of:

	LTWater
	LTGlycerol
	LTEthanol
	LTDetergent
	LTCulture
	LTProtein
	LTDNA
	LTload
	LTDoNotMix
	LTloadwater
	LTNeedToMix
	LTPostMix
	LTPreMix
	LTVISCOUS
	LTPAINT
	LTDISPENSEABOVE
	LTPEG
	LTProtoplasts
	LTCulutureReuse
	LTDNAMIX
	
The full list can be found by typing 

```bash
antharun lhpolicies
```

or looking at the liquidClass map in the following file: 


```bash
$GOPATH/src/github.com/antha-lang/antha/antha/anthalib/wtype/LiquidType.go
```

The details of these policies can be found in

```bash
$GOPATH/src/github.com/antha-lang/antha/microArch/driver/liquidhandling/makelhpolicy.go
```


## Excercises

1. Check the available plates using ```antharun lhplates``` and change inputPlateType to one of the valid alternatives in the parameters file config section

2. Check the available components and change Solution from water to one of these.

## Next Steps
Now move on to Lesson 2 where you can find out about how to perform more advanced liquid handling. 
