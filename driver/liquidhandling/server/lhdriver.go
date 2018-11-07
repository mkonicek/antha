package server

import (
	"encoding/json"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver/liquidhandling/pb"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func makeCommandReply(cs driver.CommandStatus) *pb.CommandReply {
	return &pb.CommandReply{
		OK:        cs.OK,
		Msg:       cs.Msg,
		Errorcode: cs.Errorcode,
	}
}

// liquidHandlingServer implements functionality common to low and high level servers
type liquidHandlingServer struct {
	driver liquidhandling.LiquidhandlingDriver
}

func (lhs *liquidHandlingServer) AddPlateTo(ctx context.Context, req *pb.AddPlateToRequest) (*pb.CommandReply, error) {
	if plt, err := wtype.UnmarshalLHObject([]byte(req.Plate_JSON)); err != nil {
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
	if propsJSON, err := json.Marshal(props); err != nil {
		return &pb.GetCapabilitiesReply{
			Status: makeCommandReply(cs),
		}, err
	} else {
		return &pb.GetCapabilitiesReply{
			LHProperties_JSON: propsJSON,
			Status:            makeCommandReply(cs),
		}
	}
}

func (lhs *liquidHandlingServer) GetOutputFile(context.Context, *pb.GetOutputFileRequest) (*pb.GetOutputFileReply, error) {
	b, cs := lhs.driver.GetOutputFile()
	return &ph.GetOutputFileReply{
		OutputFile: b,
		Status:     makeCommandReply(cs),
	}
}

func (lhs *liquidHandlingServer) Initialize(context.Context, *pb.InitializeRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.Initialize())
}

func (lhs *liquidHandlingServer) Message(_ context.Context, req *pb.MessageRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.Message(req.Level, req.Title, req.Text, req.ShowCancel))
}

func (lhs *liquidHandlingServer) RemoveAllPlates(context.Context, *pb.RemoveAllPlatesRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.RemoveAllPlates())
}

func (lhs *liquidHandlingServer) RemovePlateAt(_ context.Context, req *pb.RemovePlateAtRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lhs.driver.RemovePlateAt(req.Position))
}
