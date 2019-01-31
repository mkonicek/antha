package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	drv "github.com/antha-lang/antha/driver/antha_driver_v1"
	pb "github.com/antha-lang/antha/driver/liquidhandling/pb"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func toInt32(i []int) []int32 {
	r := make([]int32, 0, len(i))
	for _, v := range i {
		r = append(r, int32(v))
	}
	return r
}

func commandStatus(r *pb.CommandReply) driver.CommandStatus {
	return driver.CommandStatus{
		ErrorCode: driver.ErrorCode(r.Errorcode),
		Msg:       r.Msg,
	}
}

type LowLevelClient struct {
	client pb.LowLevelLiquidhandlingDriverClient
	driver drv.DriverClient
}

// NewLowLevelClient create a client for connecting with a remote low level
// server
func NewLowLevelClient(address string) (*LowLevelClient, error) {
	if conn, err := grpc.Dial(address, grpc.WithInsecure()); err != nil {
		return nil, errors.WithMessage(err, "Cannot initialize driver")
	} else {
		return NewLowLevelClientFromConn(conn), nil
	}
}

// NewLowLevelClientFromConn create a client for connecting with a remove low level server from a grpc Conn object
func NewLowLevelClientFromConn(conn *grpc.ClientConn) *LowLevelClient {
	return &LowLevelClient{
		client: pb.NewLowLevelLiquidhandlingDriverClient(conn),
		driver: drv.NewDriverClient(conn),
	}
}

func (llc *LowLevelClient) handleCommandReply(reply *pb.CommandReply, err error) driver.CommandStatus {
	if err != nil {
		return driver.CommandError(err.Error())
	} else {
		return commandStatus(reply)
	}
}

func (llc *LowLevelClient) DriverType() ([]string, error) {
	if reply, err := llc.driver.DriverType(context.Background(), &drv.TypeRequest{}); err != nil {
		return nil, err
	} else {
		return append([]string{reply.GetType()}, reply.GetSubtypes()...), nil
	}
}

func (llc *LowLevelClient) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	if obj, ok := plate.(wtype.LHObject); !ok {
		return driver.CommandError(fmt.Sprintf("unable to serialize object of type %T", plate))
	} else if plateJSON, err := wtype.MarshalDeckObject(obj); err != nil {
		return driver.CommandError(err.Error())
	} else {
		r, err := llc.client.AddPlateTo(context.Background(), &pb.AddPlateToRequest{
			Position:   position,
			Plate_JSON: string(plateJSON),
			Name:       name,
		})
		return llc.handleCommandReply(r, err)
	}
}

func (llc *LowLevelClient) RemoveAllPlates() driver.CommandStatus {
	r, err := llc.client.RemoveAllPlates(context.Background(), &pb.RemoveAllPlatesRequest{})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) RemovePlateAt(position string) driver.CommandStatus {
	r, err := llc.client.RemovePlateAt(context.Background(), &pb.RemovePlateAtRequest{Position: position})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Initialize() driver.CommandStatus {
	r, err := llc.client.Initialize(context.Background(), &pb.InitializeRequest{})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Finalize() driver.CommandStatus {
	r, err := llc.client.Finalize(context.Background(), &pb.FinalizeRequest{})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	r, err := llc.client.Message(context.Background(), &pb.MessageRequest{
		Level:      int32(level),
		Title:      title,
		Text:       text,
		ShowCancel: showcancel,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) GetOutputFile() ([]byte, driver.CommandStatus) {
	if r, err := llc.client.GetOutputFile(context.Background(), &pb.GetOutputFileRequest{}); err != nil {
		return nil, driver.CommandError(err.Error())
	} else {
		return r.OutputFile, commandStatus(r.Status)
	}
}

func (llc *LowLevelClient) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	if r, err := llc.client.GetCapabilities(context.Background(), &pb.GetCapabilitiesRequest{}); err != nil {
		return liquidhandling.LHProperties{}, driver.CommandError(err.Error())
	} else {
		var ret liquidhandling.LHProperties
		if err := json.Unmarshal([]byte(r.LHProperties_JSON), &ret); err != nil {
			return ret, driver.CommandError(err.Error())
		}
		return ret, commandStatus(r.Status)
	}
}

func (llc *LowLevelClient) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus {
	r, err := llc.client.Move(context.Background(), &pb.MoveRequest{
		Deckposition: deckposition,
		Wellcoords:   wellcoords,
		Reference:    toInt32(reference),
		OffsetX:      offsetX,
		OffsetY:      offsetY,
		OffsetZ:      offsetZ,
		PlateType:    plate_type,
		Head:         int32(head),
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	r, err := llc.client.Aspirate(context.Background(), &pb.AspirateRequest{
		Volume:     volume,
		Overstroke: overstroke,
		Head:       int32(head),
		Multi:      int32(multi),
		Platetype:  platetype,
		What:       what,
		Llf:        llf,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	r, err := llc.client.Dispense(context.Background(), &pb.DispenseRequest{
		Volume:    volume,
		Blowout:   blowout,
		Head:      int32(head),
		Multi:     int32(multi),
		Platetype: platetype,
		What:      what,
		Llf:       llf,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	r, err := llc.client.LoadTips(context.Background(), &pb.LoadTipsRequest{
		Channels:  toInt32(channels),
		Head:      int32(head),
		Multi:     int32(multi),
		Platetype: platetype,
		Position:  position,
		Well:      well,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	r, err := llc.client.UnloadTips(context.Background(), &pb.UnloadTipsRequest{
		Channels:  toInt32(channels),
		Head:      int32(head),
		Multi:     int32(multi),
		Platetype: platetype,
		Position:  position,
		Well:      well,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	r, err := llc.client.SetPipetteSpeed(context.Background(), &pb.SetPipetteSpeedRequest{
		Head:    int32(head),
		Channel: int32(channel),
		Rate:    rate,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	r, err := llc.client.SetDriveSpeed(context.Background(), &pb.SetDriveSpeedRequest{
		Drive: drive,
		Rate:  rate,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Wait(time float64) driver.CommandStatus {
	r, err := llc.client.Wait(context.Background(), &pb.WaitRequest{
		Time: time,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Mix(head int, volume []float64, platetype []string, cycles []int, multi int, what []string, blowout []bool) driver.CommandStatus {
	r, err := llc.client.Mix(context.Background(), &pb.MixRequest{
		Head:      int32(head),
		Volume:    volume,
		Platetype: platetype,
		Cycles:    toInt32(cycles),
		Multi:     int32(multi),
		What:      what,
		Blowout:   blowout,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) ResetPistons(head, channel int) driver.CommandStatus {
	r, err := llc.client.ResetPistons(context.Background(), &pb.ResetPistonsRequest{
		Head:    int32(head),
		Channel: int32(channel),
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
	if propsJSON, err := json.Marshal(props); err != nil {
		return driver.CommandError(err.Error())
	} else {
		r, err := llc.client.UpdateMetaData(context.Background(), &pb.UpdateMetaDataRequest{
			LHProperties_JSON: string(propsJSON),
		})
		return llc.handleCommandReply(r, err)
	}
}
