package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	grpc "google.golang.org/grpc"

	drv "github.com/antha-lang/antha/driver/antha_driver_v1"
	"github.com/antha-lang/antha/driver/liquidhandling/pb"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func toInt(s []int32) []int {
	r := make([]int, 0, len(s))
	for _, v := range s {
		r = append(r, int(v))
	}
	return r
}

// LowLevelLiquidHandlingServer a server object to listen to RPC calls to a low
// level liquid handler instruction generator
type LowLevelServer struct {
	liquidHandlingServer
	driver liquidhandling.LowLevelLiquidhandlingDriver
}

// NewLowLevelServer create a new low level server wrapping the given driver
func NewLowLevelServer(driver liquidhandling.LowLevelLiquidhandlingDriver) (*LowLevelServer, error) {
	return &LowLevelServer{
		liquidHandlingServer: liquidHandlingServer{
			driver: driver,
		},
		driver: driver,
	}, nil
}

// Listen begin listening for gRPC calls on the given port. returns only in error
func (lhs *LowLevelServer) Listen(port int) error {
	if lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	} else {
		fmt.Println("Listening at", lis.Addr().String())

		s := grpc.NewServer()
		pb.RegisterLowLevelLiquidhandlingDriverServer(s, lhs)
		drv.RegisterDriverServer(s, lhs)
		return s.Serve(lis)
	}
}

func (lls *LowLevelServer) Aspirate(_ context.Context, req *pb.AspirateRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Aspirate(req.Volume, req.Overstroke, int(req.Head), int(req.Multi), req.Platetype, req.What, req.Llf)), nil
}

func (lls *LowLevelServer) Dispense(_ context.Context, req *pb.DispenseRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Dispense(req.Volume, req.Blowout, int(req.Head), int(req.Multi), req.Platetype, req.What, req.Llf)), nil
}

func (lls *LowLevelServer) LoadTips(_ context.Context, req *pb.LoadTipsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.LoadTips(toInt(req.Channels), int(req.Head), int(req.Multi), req.Platetype, req.Position, req.Well)), nil
}

func (lls *LowLevelServer) Mix(_ context.Context, req *pb.MixRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Mix(int(req.Head), req.Volume, req.Platetype, toInt(req.Cycles), int(req.Multi), req.What, req.Blowout)), nil
}

func (lls *LowLevelServer) Move(_ context.Context, req *pb.MoveRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Move(req.Deckposition, req.Wellcoords, toInt(req.Reference), req.OffsetX, req.OffsetY, req.OffsetZ, req.PlateType, int(req.Head))), nil
}

func (lls *LowLevelServer) ResetPistons(_ context.Context, req *pb.ResetPistonsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.ResetPistons(int(req.Head), int(req.Channel))), nil
}

func (lls *LowLevelServer) SetDriveSpeed(_ context.Context, req *pb.SetDriveSpeedRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.SetDriveSpeed(req.Drive, req.Rate)), nil
}

func (lls *LowLevelServer) SetPipetteSpeed(_ context.Context, req *pb.SetPipetteSpeedRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.SetPipetteSpeed(int(req.Head), int(req.Channel), req.Rate)), nil
}

func (lls *LowLevelServer) UnloadTips(_ context.Context, req *pb.UnloadTipsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.UnloadTips(toInt(req.Channels), int(req.Head), int(req.Multi), req.Platetype, req.Position, req.Well)), nil
}

func (lls *LowLevelServer) UpdateMetaData(_ context.Context, req *pb.UpdateMetaDataRequest) (*pb.CommandReply, error) {
	var props liquidhandling.LHProperties
	if err := json.Unmarshal([]byte(req.LHProperties_JSON), &props); err != nil {
		return &pb.CommandReply{Msg: err.Error()}, err
	} else {
		return makeCommandReply(lls.driver.UpdateMetaData(&props)), nil
	}
}

func (lls *LowLevelServer) Wait(_ context.Context, req *pb.WaitRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Wait(req.Time)), nil
}
