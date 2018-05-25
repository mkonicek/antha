package lh

import (
	"fmt"

	pb "github.com/antha-lang/antha/driver/pb/lh"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type LLLHDriver struct {
	C pb.LowLevelLiquidhandlingDriverClient
}

func NewLLLHDriver(address string) *LLLHDriver {
	var d LLLHDriver
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize driver: %s", err))
	}

	d.C = pb.NewLowLevelLiquidhandlingDriverClient(conn)

	return &d
}

func (d *LLLHDriver) AddPlateTo(arg_1 string, arg_2 interface{}, arg_3 string) driver.CommandStatus {
	req := pb.AddPlateToRequest{
		(string)(arg_1),
		Encodeinterface(arg_2),
		(string)(arg_3),
	}
	ret, _ := d.C.AddPlateTo(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) Finalize() driver.CommandStatus {
	req := pb.FinalizeRequest{}
	ret, _ := d.C.Finalize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	req := pb.GetCapabilitiesRequest{}
	ret, err := d.C.GetCapabilities(context.Background(), &req)
	if err != nil {
		return liquidhandling.LHProperties{}, driver.CommandStatus{
			Msg: err.Error(),
		}
	}
	return (liquidhandling.LHProperties)(DecodeLHProperties(ret.Ret_1)), (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}

func (d *LLLHDriver) GetOutputFile() ([]byte, driver.CommandStatus) {
	req := pb.GetOutputFileRequest{}

	ret, _ := d.C.GetOutputFile(context.Background(), &req)
	return ret.Ret_1, (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *LLLHDriver) Initialize() driver.CommandStatus {
	req := pb.InitializeRequest{}
	ret, _ := d.C.Initialize(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) Message(arg_1 int, arg_2 string, arg_3 string, arg_4 bool) driver.CommandStatus {
	req := pb.MessageRequest{
		int64(arg_1),
		(string)(arg_2),
		(string)(arg_3),
		(bool)(arg_4),
	}
	ret, _ := d.C.Message(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) RemoveAllPlates() driver.CommandStatus {
	req := pb.RemoveAllPlatesRequest{}
	ret, _ := d.C.RemoveAllPlates(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) RemovePlateAt(arg_1 string) driver.CommandStatus {
	req := pb.RemovePlateAtRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.RemovePlateAt(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}

func (d *LLLHDriver) Aspirate(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) driver.CommandStatus {
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
func (d *LLLHDriver) Dispense(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) driver.CommandStatus {
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
func (d *LLLHDriver) LoadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) driver.CommandStatus {
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
func (d *LLLHDriver) Mix(arg_1 int, arg_2 []float64, arg_3 []string, arg_4 []int, arg_5 int, arg_6 []string, arg_7 []bool) driver.CommandStatus {
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
func (d *LLLHDriver) Move(arg_1 []string, arg_2 []string, arg_3 []int, arg_4 []float64, arg_5 []float64, arg_6 []float64, arg_7 []string, arg_8 int) driver.CommandStatus {
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
func (d *LLLHDriver) ResetPistons(arg_1 int, arg_2 int) driver.CommandStatus {
	req := pb.ResetPistonsRequest{
		int64(arg_1),
		int64(arg_2),
	}
	ret, _ := d.C.ResetPistons(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) SetDriveSpeed(arg_1 string, arg_2 float64) driver.CommandStatus {
	req := pb.SetDriveSpeedRequest{
		(string)(arg_1),
		(float64)(arg_2),
	}
	ret, _ := d.C.SetDriveSpeed(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) SetPipetteSpeed(arg_1 int, arg_2 int, arg_3 float64) driver.CommandStatus {
	req := pb.SetPipetteSpeedRequest{
		int64(arg_1),
		int64(arg_2),
		(float64)(arg_3),
	}
	ret, _ := d.C.SetPipetteSpeed(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *LLLHDriver) UnloadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) driver.CommandStatus {
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
func (d *LLLHDriver) UpdateMetaData(arg_1 *liquidhandling.LHProperties) driver.CommandStatus {
	req := pb.UpdateMetaDataRequest{
		EncodePtrToLHProperties(arg_1),
	}
	ret, _ := d.C.UpdateMetaData(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}

/*
func (d *LLLHDriver) MoveRaw(arg_1 int, arg_2 float64, arg_3 float64, arg_4 float64) driver.CommandStatus {
	req := pb.MoveRawRequest{
		int64(arg_1),
		(float64)(arg_2),
		(float64)(arg_3),
		(float64)(arg_4),
	}
	ret, _ := d.C.MoveRaw(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
*/

func (d *LLLHDriver) Wait(arg_1 float64) driver.CommandStatus {
	req := pb.WaitRequest{
		(float64)(arg_1),
	}
	ret, _ := d.C.Wait(context.Background(), &req)
	return (driver.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
