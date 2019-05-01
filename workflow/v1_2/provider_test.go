package v1_2

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func getTestV1_2WorkflowProvider() (*V1_2WorkflowProvider, error) {
	fixture := filepath.Join("testdata", "sample_v1_2_workflow.json")
	bytes, err := ioutil.ReadFile(fixture)
	if err != nil {
		return nil, err
	}

	wf := &workflowv1_2{}
	err = json.Unmarshal(bytes, wf)
	if err != nil {
		return nil, err
	}

	p := NewV1_2WorkflowProvider(wf)

	return p, nil
}

func TestGetMeta(t *testing.T) {
	p, err := getTestV1_2WorkflowProvider()
	if err != nil {
		t.Fatal(err)
	}

	m, err := p.GetMeta()
	if err != nil {
		t.Fatal(err)
	}

	if m == nil {
		t.Fatal("Got nil Meta from GetMeta()")
	}

	expectedName := "My Test Workflow"
	if m.Name != expectedName {
		t.Errorf("Expected name '%v', got '%v'", expectedName, m.Name)
	}
}
