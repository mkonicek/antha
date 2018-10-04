package lh

import (
	"fmt"

	pb "github.com/antha-lang/antha/driver/pb/lh"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type HLLHDriver struct {
	C pb.HighLevelLiquidhandlingDriverClient
}

func NewHLLHDriver(address string) *HLLHDriver {
	var d HLLHDriver
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize driver: %s", err))
	}

	d.C = pb.NewHighLevelLiquidhandlingDriverClient(conn)

	return &d
}

func (d *HLLHDriver) AddPlateTo(arg_1 string, arg_2 interface{}, arg_3 string) driver.CommandStatus {
	req := pb.AddPlateToRequest{
		(string)(arg_1),
		Encodeinterface(arg_2),
		(string)(arg_3),
	}
	ret, _ := d.C.AddPlateTo(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *HLLHDriver) Finalize() driver.CommandStatus {
	req := pb.FinalizeRequest{}
	ret, _ := d.C.Finalize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *HLLHDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	req := pb.GetCapabilitiesRequest{}
	ret, err := d.C.GetCapabilities(context.Background(), &req)
	if err != nil {
		return liquidhandling.LHProperties{}, driver.CommandStatus{
			Msg: err.Error(),
		}
	}
	return (liquidhandling.LHProperties)(DecodeLHProperties(ret.Ret_1)), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}

func (d *HLLHDriver) GetOutputFile() ([]byte, driver.CommandStatus) {
	req := pb.GetOutputFileRequest{}

	ret, _ := d.C.GetOutputFile(context.Background(), &req)
	return ret.Ret_1, (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *HLLHDriver) Initialize() driver.CommandStatus {
	req := pb.InitializeRequest{}
	ret, _ := d.C.Initialize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *HLLHDriver) Message(arg_1 int, arg_2 string, arg_3 string, arg_4 bool) driver.CommandStatus {
	req := pb.MessageRequest{
		int64(arg_1),
		(string)(arg_2),
		(string)(arg_3),
		(bool)(arg_4),
	}
	ret, _ := d.C.Message(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *HLLHDriver) RemoveAllPlates() driver.CommandStatus {
	req := pb.RemoveAllPlatesRequest{}
	ret, _ := d.C.RemoveAllPlates(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *HLLHDriver) RemovePlateAt(arg_1 string) driver.CommandStatus {
	req := pb.RemovePlateAtRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.RemovePlateAt(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}

func (d *HLLHDriver) Transfer(what, platefrom, wellfrom, plateto, wellto []string, volume []float64) driver.CommandStatus {
	req := pb.TransferRequest{
		what, platefrom, wellfrom, plateto, wellto, volume,
	}

	ret, _ := d.C.Transfer(context.Background(), &req)

	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
