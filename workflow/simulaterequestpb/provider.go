package simulaterequestpb

import (
	"io"

	"github.com/Synthace/antha-runner/protobuf"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/golang/protobuf/proto"
)

type SimulateRequestProtobufProvider struct {
	pb               *protobuf.SimulateRequest
	fm               *effects.FileManager
	repoMap          workflow.ElementTypesByRepository
	gilsonDeviceName string
	logger           *logger.Logger
}

func NewProvider(
	oldWorkflowReader io.Reader,
	fm *effects.FileManager,
	repoMap workflow.ElementTypesByRepository,
	gilsonDeviceName string,
	logger *logger.Logger,
) (*SimulateRequestProtobufProvider, error) {
	var bytes []byte
	if _, err := oldWorkflowReader.Read(bytes); err != nil {
		return nil, err
	}

	pb := &protobuf.SimulateRequest{}
	if err := proto.Unmarshal(bytes, pb); err != nil {
		return nil, err
	}

	return &SimulateRequestProtobufProvider{
		pb:               pb,
		fm:               fm,
		repoMap:          repoMap,
		gilsonDeviceName: gilsonDeviceName,
		logger:           logger,
	}, nil

	return nil, nil
}

func (p *SimulateRequestProtobufProvider) GetWorkflowID() (workflow.BasicId, error) {
	return "", nil
}

func (p *SimulateRequestProtobufProvider) GetMeta() (workflow.Meta, error) {
	return workflow.Meta{}, nil
}

func (p *SimulateRequestProtobufProvider) GetRepositories() (workflow.Repositories, error) {
	return workflow.Repositories{}, nil
}

func (p *SimulateRequestProtobufProvider) GetElements() (workflow.Elements, error) {
	return workflow.Elements{}, nil
}

func (p *SimulateRequestProtobufProvider) GetInventory() (workflow.Inventory, error) {
	return workflow.Inventory{}, nil
}

func (p *SimulateRequestProtobufProvider) GetConfig() (workflow.Config, error) {
	return workflow.Config{}, nil
}

func (p *SimulateRequestProtobufProvider) GetTesting() (workflow.Testing, error) {
	return workflow.Testing{}, nil
}
