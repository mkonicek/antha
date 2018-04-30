package lh

import (
	"fmt"

	pb "github.com/antha-lang/antha/driver/pb/lh"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Driver struct {
	C pb.ExtendedLiquidhandlingDriverClient
}

func NewDriver(address string) *Driver {
	var d Driver
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize driver: %s", err))
	}

	d.C = pb.NewExtendedLiquidhandlingDriverClient(conn)

	return &d
}

func (d *Driver) AddPlateTo(arg_1 string, arg_2 interface{}, arg_3 string) driver.CommandStatus {
	req := pb.AddPlateToRequest{
		(string)(arg_1),
		Encodeinterface(arg_2),
		(string)(arg_3),
	}
	ret, _ := d.C.AddPlateTo(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Aspirate(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) driver.CommandStatus {
	req := pb.AspirateRequest{
		EncodeArrayOffloat64(arg_1),
		EncodeArrayOfbool(arg_2),
		int64(arg_3),
		int64(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
		EncodeArrayOfbool(arg_7),
	}
	ret, _ := d.C.Aspirate(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Close() driver.CommandStatus {
	req := pb.CloseRequest{}
	ret, _ := d.C.Close(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Dispense(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) driver.CommandStatus {
	req := pb.DispenseRequest{
		EncodeArrayOffloat64(arg_1),
		EncodeArrayOfbool(arg_2),
		int64(arg_3),
		int64(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
		EncodeArrayOfbool(arg_7),
	}
	ret, _ := d.C.Dispense(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Finalize() driver.CommandStatus {
	req := pb.FinalizeRequest{}
	ret, _ := d.C.Finalize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	req := pb.GetCapabilitiesRequest{}
	ret, err := d.C.GetCapabilities(context.Background(), &req)
	if err != nil {
		return liquidhandling.LHProperties{}, driver.CommandStatus{
			Msg: err.Error(),
		}
	}
	return (liquidhandling.LHProperties)(DecodeLHProperties(ret.Ret_1)), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetCurrentPosition(arg_1 int) (string, driver.CommandStatus) {
	req := pb.GetCurrentPositionRequest{
		int64(arg_1),
	}
	ret, _ := d.C.GetCurrentPosition(context.Background(), &req)
	return (string)(ret.Ret_1), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}

func (d *Driver) GetOutputFile() (string, driver.CommandStatus) {
	req := pb.GetOutputFileRequest{}

	ret, _ := d.C.GetOutputFile(context.Background(), &req)
	return (string)(ret.Ret_1), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetHeadState(arg_1 int) (string, driver.CommandStatus) {
	req := pb.GetHeadStateRequest{
		int64(arg_1),
	}
	ret, _ := d.C.GetHeadState(context.Background(), &req)
	return (string)(ret.Ret_1), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetPositionState(arg_1 string) (string, driver.CommandStatus) {
	req := pb.GetPositionStateRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.GetPositionState(context.Background(), &req)
	return (string)(ret.Ret_1), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetStatus() (driver.Status, driver.CommandStatus) {
	req := pb.GetStatusRequest{}
	ret, _ := d.C.GetStatus(context.Background(), &req)
	return (driver.Status)(DecodeMapstringinterfaceMessage(ret.Ret_1)), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) Go() driver.CommandStatus {
	req := pb.GoRequest{}
	ret, _ := d.C.Go(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Initialize() driver.CommandStatus {
	req := pb.InitializeRequest{}
	ret, _ := d.C.Initialize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LightsOff() driver.CommandStatus {
	req := pb.LightsOffRequest{}
	ret, _ := d.C.LightsOff(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LightsOn() driver.CommandStatus {
	req := pb.LightsOnRequest{}
	ret, _ := d.C.LightsOn(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadAdaptor(arg_1 int) driver.CommandStatus {
	req := pb.LoadAdaptorRequest{
		int64(arg_1),
	}
	ret, _ := d.C.LoadAdaptor(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadHead(arg_1 int) driver.CommandStatus {
	req := pb.LoadHeadRequest{
		int64(arg_1),
	}
	ret, _ := d.C.LoadHead(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) driver.CommandStatus {
	req := pb.LoadTipsRequest{
		EncodeArrayOfint(arg_1),
		int64(arg_2),
		int64(arg_3),
		EncodeArrayOfstring(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
	}
	ret, _ := d.C.LoadTips(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Message(arg_1 int, arg_2 string, arg_3 string, arg_4 bool) driver.CommandStatus {
	req := pb.MessageRequest{
		int64(arg_1),
		(string)(arg_2),
		(string)(arg_3),
		(bool)(arg_4),
	}
	ret, _ := d.C.Message(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Mix(arg_1 int, arg_2 []float64, arg_3 []string, arg_4 []int, arg_5 int, arg_6 []string, arg_7 []bool) driver.CommandStatus {
	req := pb.MixRequest{
		int64(arg_1),
		EncodeArrayOffloat64(arg_2),
		EncodeArrayOfstring(arg_3),
		EncodeArrayOfint(arg_4),
		int64(arg_5),
		EncodeArrayOfstring(arg_6),
		EncodeArrayOfbool(arg_7),
	}
	ret, _ := d.C.Mix(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Move(arg_1 []string, arg_2 []string, arg_3 []int, arg_4 []float64, arg_5 []float64, arg_6 []float64, arg_7 []string, arg_8 int) driver.CommandStatus {
	req := pb.MoveRequest{
		EncodeArrayOfstring(arg_1),
		EncodeArrayOfstring(arg_2),
		EncodeArrayOfint(arg_3),
		EncodeArrayOffloat64(arg_4),
		EncodeArrayOffloat64(arg_5),
		EncodeArrayOffloat64(arg_6),
		EncodeArrayOfstring(arg_7),
		int64(arg_8),
	}
	ret, _ := d.C.Move(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) MoveRaw(arg_1 int, arg_2 float64, arg_3 float64, arg_4 float64) driver.CommandStatus {
	req := pb.MoveRawRequest{
		int64(arg_1),
		(float64)(arg_2),
		(float64)(arg_3),
		(float64)(arg_4),
	}
	ret, _ := d.C.MoveRaw(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Open() driver.CommandStatus {
	req := pb.OpenRequest{}
	ret, _ := d.C.Open(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) RemoveAllPlates() driver.CommandStatus {
	req := pb.RemoveAllPlatesRequest{}
	ret, _ := d.C.RemoveAllPlates(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) RemovePlateAt(arg_1 string) driver.CommandStatus {
	req := pb.RemovePlateAtRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.RemovePlateAt(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) ResetPistons(arg_1 int, arg_2 int) driver.CommandStatus {
	req := pb.ResetPistonsRequest{
		int64(arg_1),
		int64(arg_2),
	}
	ret, _ := d.C.ResetPistons(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetDriveSpeed(arg_1 string, arg_2 float64) driver.CommandStatus {
	req := pb.SetDriveSpeedRequest{
		(string)(arg_1),
		(float64)(arg_2),
	}
	ret, _ := d.C.SetDriveSpeed(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetPipetteSpeed(arg_1 int, arg_2 int, arg_3 float64) driver.CommandStatus {
	req := pb.SetPipetteSpeedRequest{
		int64(arg_1),
		int64(arg_2),
		(float64)(arg_3),
	}
	ret, _ := d.C.SetPipetteSpeed(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetPositionState(arg_1 string, arg_2 driver.PositionState) driver.CommandStatus {
	req := pb.SetPositionStateRequest{
		(string)(arg_1),
		(EncodeMapstringinterfaceMessage(arg_2)),
	}
	ret, _ := d.C.SetPositionState(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Stop() driver.CommandStatus {
	req := pb.StopRequest{}
	ret, _ := d.C.Stop(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadAdaptor(arg_1 int) driver.CommandStatus {
	req := pb.UnloadAdaptorRequest{
		int64(arg_1),
	}
	ret, _ := d.C.UnloadAdaptor(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadHead(arg_1 int) driver.CommandStatus {
	req := pb.UnloadHeadRequest{
		int64(arg_1),
	}
	ret, _ := d.C.UnloadHead(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) driver.CommandStatus {
	req := pb.UnloadTipsRequest{
		EncodeArrayOfint(arg_1),
		int64(arg_2),
		int64(arg_3),
		EncodeArrayOfstring(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
	}
	ret, _ := d.C.UnloadTips(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UpdateMetaData(arg_1 *liquidhandling.LHProperties) driver.CommandStatus {
	req := pb.UpdateMetaDataRequest{
		EncodePtrToLHProperties(arg_1),
	}
	ret, _ := d.C.UpdateMetaData(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Wait(arg_1 float64) driver.CommandStatus {
	req := pb.WaitRequest{
		(float64)(arg_1),
	}
	ret, _ := d.C.Wait(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
