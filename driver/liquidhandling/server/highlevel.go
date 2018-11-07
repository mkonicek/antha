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

// HighLevelLiquidHandlingServer a server object to listen to RPC calls to a high
// level liquid handler instruction generator
type HighLevelServer struct {
	liquidHandlingServer
	driver liquidhandling.HighLevelLiquidhandlingDriver
}

// NewHighLevelServer create a new high level server wrapping the given driver
func NewHighLevelServer(driver liquidhandling.HighLevelLiquidHandlingDriver) (*HighLevelServer, error) {
	return &HighLevelServer{
		liquidHandlingServer: &liquidHandlingServer{
			driver: driver,
		},
		driver: driver,
	}
}

// Listen begin listening for gRPC calls on the given port. returns only in error
func (lhs *HighLevelServer) Listen(port int) error {
	if lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	} else {
		fmt.Println("Listening at", lis.Addr().String())

		s := grpc.NewServer()
		pb.RegisterHighLevelLiquidhandlingDriverServer(s, lhs)
		s.Serve(lis)
	}
}

func (lhs *HighLevelServer) Transfer(_ context.Context, req *pb.TransferRequest) (*pb.CommandReply, error) {
	return makeCommandReply(hls.driver.Transfer(req.What, req.PlateFrom, req.Wellfrom, req.Plateto, req.Wellto, req.Volume))
}
