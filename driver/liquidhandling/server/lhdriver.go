package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	drv "github.com/antha-lang/antha/driver/antha_driver_v1"
	"github.com/antha-lang/antha/driver/liquidhandling/pb"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func makeCommandReply(cs driver.CommandStatus) *pb.CommandReply {
	return &pb.CommandReply{
		Msg:       cs.Msg,
		Errorcode: int32(cs.ErrorCode),
	}
}

// liquidHandlingServer implements functionality common to low and high level servers
type liquidHandlingServer struct {
	driver liquidhandling.LiquidhandlingDriver
}

func (lhs *liquidHandlingServer) DriverType(_ context.Context, req *drv.TypeRequest) (*drv.TypeReply, error) {
	if strs, err := lhs.driver.DriverType(); err != nil {
		return nil, err
	} else if len(strs) == 0 {
		return nil, errors.New("DriverType() returned no values")
	} else {
		return &drv.TypeReply{Type: strs[0], Subtypes: strs[1:]}, nil
	}
}

func (lhs *liquidHandlingServer) AddPlateTo(ctx context.Context, req *pb.AddPlateToRequest) (*pb.CommandReply, error) {
	if plt, err := wtype.UnmarshalDeckObject([]byte(req.Plate_JSON)); err != nil {
		return nil, err
	} else {
		return makeCommandReply(lhs.driver.AddPlateTo(req.Position, plt, req.Name)), nil
	}
}

func (lhs *liquidHandlingServer) Finalize(context.Context, *pb.FinalizeRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.Finalize()), nil
}

func (lhs *liquidHandlingServer) GetCapabilities(context.Context, *pb.GetCapabilitiesRequest) (*pb.GetCapabilitiesReply, error) {
	props, cs := lhs.driver.GetCapabilities()

	if propsJSON, err := json.Marshal(&props); err != nil {
		return &pb.GetCapabilitiesReply{
			Status: makeCommandReply(cs),
		}, err
	} else {
		return &pb.GetCapabilitiesReply{
			LHProperties_JSON: string(propsJSON),
			Status:            makeCommandReply(cs),
		}, nil
	}
}

func (lhs *liquidHandlingServer) GetOutputFile(context.Context, *pb.GetOutputFileRequest) (*pb.GetOutputFileReply, error) {
	b, cs := lhs.driver.GetOutputFile()
	return &pb.GetOutputFileReply{
		OutputFile: b,
		Status:     makeCommandReply(cs),
	}, nil
}

func (lhs *liquidHandlingServer) Initialize(context.Context, *pb.InitializeRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.Initialize()), nil
}

func (lhs *liquidHandlingServer) Message(_ context.Context, req *pb.MessageRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.Message(int(req.Level), req.Title, req.Text, req.ShowCancel)), nil
}

func (lhs *liquidHandlingServer) RemoveAllPlates(context.Context, *pb.RemoveAllPlatesRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.RemoveAllPlates()), nil
}

func (lhs *liquidHandlingServer) RemovePlateAt(_ context.Context, req *pb.RemovePlateAtRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.RemovePlateAt(req.Position)), nil
}
