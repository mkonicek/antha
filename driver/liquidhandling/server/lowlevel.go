package server

import (
	"encoding/json"
	"fmt"
	grpc "google.golang.org/grpc"
	"net"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver/liquidhandling/pb"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// LowLevelLiquidHandlingServer a server object to listen to RPC calls to a low
// level liquid handler instruction generator
type LowLevelServer struct {
	liquidHandlingServer
	driver liquidhandling.LowLevelLiquidhandlingDriver
}

// NewLowLevelServer create a new low level server wrapping the given driver
func NewLowLevelServer(driver liquidhandling.LowLevelLiquidHandlingDriver) (*LowLevelServer, error) {
	return &LowLevelServer{
		liquidHandlingServer: &liquidHandlingServer{
			driver: driver,
		},
		driver: driver,
	}
}

// Listen begin listening for gRPC calls on the given port. returns only in error
func (lhs *LowLevelServer) Listen(port int) error {
	if lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	} else {
		fmt.Println("Listening at", lis.Addr().String())

		s := grpc.NewServer()
		pb.RegisterLowLevelLiquidhandlingDriverServer(s, lhs)
		s.Serve(lis)
	}
}

func (lls *LowLevelServer) Aspirate(_ context.Context, req *pb.AspirateRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Aspirate(req.Volume, req.Overstroke, req.Head, req.Multi, req.Platetype, req.What, req.Llf)), nil
}

func (lls *LowLevelServer) Dispense(_ context.Context, req *pb.DispenseRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Dispense(req.Volume, req.Blowout, req.Head, req.Multi, req.Platetype, req.What, req.Llf)), nil
}

func (lls *LowLevelServer) LoadTips(_ context.Context, req *pb.LoadTipsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.LoadTips(req.Channels, req.Head, req.Multi, req.Platetype, req.Position, req.Well)), nil
}

func (lls *LowLevelServer) Mix(_ context.Context, req *pb.MixRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Mix(req.Head, req.Volume, req.Platetype, req.Cycles, req.Multi, req.What, req.Blowout)), nil
}

func (lls *LowLevelServer) Move(_ context.Context, req *pb.MoveRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Move(req.Deckposition, req.Wellcoords, req.Reference, req.OffsetX, req.OffsetY, req.OffsetZ, req.Platetype, req.Head)), nil
}

func (lls *LowLevelServer) ResetPistons(_ context.Context, req *pb.ResetPistonsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.ResetPistons(req.Head, req.Channels)), nil
}

func (lls *LowLevelServer) SetDriveSpeed(_ context.Context, req *pb.SetDriveSpeedRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.SetDriveSpeed(req.Drive, req.Rate)), nil
}

func (lls *LowLevelServer) SetPipetteSpeed(_ context.Context, req *pb.SetPipetteSpeedRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.SetPipetteSpeed(req.Head, req.Channel, req.Rate)), nil
}

func (lls *LowLevelServer) UnloadTips(_ context.Context, req *pb.UnloadTipsRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.LoadTips(req.Channels, req.Head, req.Multi, req.Platetype, req.Position, req.Well)), nil
}

func (lls *LowLevelServer) Wait(_ context.Context, req *pb.WaitRequest) (*pb.CommandReply, error) {
	return makeCommandReply(lls.driver.Wait(req.Time)), nil
}
