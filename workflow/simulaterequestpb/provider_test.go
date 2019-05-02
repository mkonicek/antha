package simulaterequestpb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflow/migrate"
	"github.com/antha-lang/antha/workflow/simulaterequestpb"
)

func getTestProvider() (migrate.WorkflowProvider, error) {
	protobufFilePath := filepath.Join("testdata", "request.pb")
	tmpDir, err := ioutil.TempDir("", "tests")
	if err != nil {
		return nil, err
	}
	fm, err := effects.NewFileManager(tmpDir, tmpDir)
	if err != nil {
		return nil, err
	}

	repo := &workflow.Repository{
		Directory: "/tmp",
	}

	elementNames := []string{"AccuracyTest", "Aliquot_Liquid"}

	elementTypeMap := workflow.ElementTypeMap{}
	for _, name := range elementNames {
		etn := workflow.ElementTypeName(name)
		ep := workflow.ElementPath("Elements/Test/" + name)
		elementTypeMap[etn] = workflow.ElementType{
			RepositoryName: "repos.antha.com/antha-test/elements-test",
			ElementPath:    ep,
		}
	}
	repoMap := workflow.ElementTypesByRepository{}
	repoMap[repo] = elementTypeMap

	gilsonDeviceName := "testie"

	logger := logger.NewLogger()

	r, err := os.Open(protobufFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return simulaterequestpb.NewProvider(r, fm, repoMap, gilsonDeviceName, logger)
}

func TestGetMeta(t *testing.T) {
	p, err := getTestProvider()
	if err != nil {
		t.Fatal(err)
	}

	m, err := p.GetMeta()
	if err != nil {
		t.Fatal(err)
	}

	expectedName := "My Test Workflow"
	if m.Name != expectedName {
		t.Errorf("Expected name '%v', got '%v'", expectedName, m.Name)
	}
}
