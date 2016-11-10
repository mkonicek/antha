package param

import (
	"strings"
	"testing"
)

func TestNoExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`{"Parameters": { "TheAlpha": { "Key": "Value" } } }`)
	if _, _, err := TryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	}
}

func TestExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }, 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }
	]`)
	if desc, params, err := TryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	} else if l := len(desc.Processes); l != 2 {
		t.Errorf("expecting workflow of %d processes but found only %d", 2, l)
	} else if l := len(params.Parameters); l != 2 {
		t.Errorf("expecting parameters for %d processes but found only %d", 2, l)
	}
}

func TestBadExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{"Parameters": { "TheAlpha": { "Key": "Value" } }, "Config": { "A": "B" } }, 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }
	]`)
	if _, p, err := TryExpand(wdata, pdata); err == nil {
		t.Errorf("expecting error but got success, output %+v", p)
	}
}

func TestInvalidInput(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha1": { "component": "Alpha" }, "TheAlpha2": { "component": "Alpha" } } }`)
	pdataOk := []byte(`{
		"Parameters": {
			"TheAlpha1": { "Key": "Value", "invalid": false },
			"TheAlpha2": { "Key": "Value" } 
		} 
	}`)
	pdataBad := []byte(strings.Replace(string(pdataOk), `"invalid":`, `invalid}}}},`, 1))

	if _, _, err := TryExpand(wdata, pdataOk); err != nil {
		t.Error(err)
	} else if _, p, err := TryExpand(wdata, pdataBad); err == nil {
		t.Errorf("expecting error but got success, output: %+v", p)
	}
}

func TestSExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{ "TheAlpha": { "Key": "Value" } }, 
		{ "TheAlpha": { "Key": "Value" } }
	]`)
	if desc, params, err := TryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	} else if l := len(desc.Processes); l != 2 {
		t.Errorf("expecting workflow of %d processes but found only %d", 2, l)
	} else if l := len(params.Parameters); l != 2 {
		t.Errorf("expecting parameters for %d processes but found only %d", 2, l)
	}
}
