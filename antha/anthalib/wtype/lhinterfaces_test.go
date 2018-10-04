package wtype

import (
	"testing"
)

func getLHObjects() map[string]interface{} {
	ret := map[string]interface{}{
		"LHComponent": &Liquid{},
		"LHDeck":      &LHDeck{},
		"LHPlate":     &Plate{},
		"LHTip":       &LHTip{},
		"LHTipbox":    &LHTipbox{},
		"LHTipwaste":  &LHTipwaste{},
		"LHWell":      &LHWell{},
	}
	return ret
}

func getAssertionMap() map[string]func(interface{}) bool {
	ret := map[string]func(interface{}) bool{
		"Named": func(obj interface{}) bool {
			_, ok := obj.(Named)
			return ok
		},
		"Identifiable": func(obj interface{}) bool {
			_, ok := obj.(Identifiable)
			return ok
		},
		"Typed": func(obj interface{}) bool {
			_, ok := obj.(Typed)
			return ok
		},
		"Classy": func(obj interface{}) bool {
			_, ok := obj.(Classy)
			return ok
		},
		"LHObject": func(obj interface{}) bool {
			_, ok := obj.(LHObject)
			return ok
		},
		"LHParent": func(obj interface{}) bool {
			_, ok := obj.(LHParent)
			return ok
		},
		"Addressable": func(obj interface{}) bool {
			_, ok := obj.(Addressable)
			return ok
		},
		"LHContainer": func(obj interface{}) bool {
			_, ok := obj.(LHContainer)
			return ok
		},
	}
	return ret
}

func TestInterfaceImplementations(t *testing.T) {

	tests := map[string][]string{
		"LHComponent": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			//"LHObject",
			//"LHParent",
			//"Addressable",
			//"LHContainer",
		},
		"LHDeck": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			"LHParent",
			//"Addressable",
			//"LHContainer",
		},
		"LHPlate": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			//"LHParent",
			"Addressable",
			//"LHContainer",
		},
		"LHTip": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			//"LHParent",
			//"Addressable",
			"LHContainer",
		},
		"LHTipbox": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			//"LHParent",
			"Addressable",
			//"LHContainer",
		},
		"LHTipwaste": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			//"LHParent",
			"Addressable",
			//"LHContainer",
		},
		"LHWell": {
			"Named",
			"Identifiable",
			"Typed",
			"Classy",
			"LHObject",
			//"LHParent",
			//"Addressable",
			"LHContainer",
		},
	}

	objects := getLHObjects()
	asserts := getAssertionMap()

	for type_name, interfaces := range tests {
		obj := objects[type_name]
		for _, interface_name := range interfaces {
			if !asserts[interface_name](obj) {
				t.Errorf("%s doesn't implement %s", type_name, interface_name)
			}
		}
	}
}
