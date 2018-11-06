package client

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/antha-lang/antha/driver/liquidhandling/pb"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type LowLevelClient struct {
	client pb.LowLevelLiquidhandlingDriverClient
}

// NewLowLevelClient create a client for connecting with a remote high level
// server
func NewLowLevelClient(address string) (*LowLevelClient, error) {
	if conn, err := grpc.Dial(address, grpc.WithInsecure()); err != nil {
		return nil, errors.WithMessage(err, "Cannot initialize driver")
	} else {
		return &LowLevelClient{
			client: pb.NewLowLevelLiquidhandlingDriverClient(conn),
		}
	}
}

func (llc *LowLevelClient) handleCommandReply(reply *pb.CommandReply, err error) driver.CommandStatus {
	if err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: r.Errorcode, Msg: r.Msg}
	}
}

func (llc *LowLevelClient) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	if plateJSON, err := json.Marshal(plate); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		r, err := llc.client.AddPlateTo(context.Background(), &pb.AddPlateToRequest{
			Position: position,
			Plate:    string(plateJSON),
			Name:     name,
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
		Level:      level,
		Title:      title,
		Text:       text,
		ShowCancel: showcancel,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) GetOutputFile() ([]byte, driver.CommandStatus) {
	if r, err := llc.client.GetOutputFile(context.Background(), &pb.GetOutputFileRequest{}); err != nil {
		return nil, driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return r.OutputFile, driver.CommandStatus{OK: r.Status.OK, Errorcode: r.Status.Errorcode, Msg: r.Status.Msg}
	}
}

func (llc *LowLevelClient) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	if r, err := llc.client.GetCapabilities(context.Background(), &pb.GetCapabilities{}); err != nil {
		return nil, driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		var ret liquidhandling.LHProperties
		if err := json.Unmarshal([]byte(r.LHProperties_JSON), &ret); err != nil {
			return ret, driver.CommandStatus{Errorcode: driver.ERR, Msg: err.Error()}
		}
		return ret, driver.CommandStatus{OK: r.Status.OK, Errorcode: r.Status.Errorcode, Msg: r.Status.Msg}
	}
}

func (llc *LowLevelClient) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus {
	r, err := llc.client.Move(context.Background(), &pb.MoveRequest{
		Deckposition: deckposition,
		Wellcoords:   wellcoords,
		Reference:    reference,
		OffsetX:      offsetX,
		OffsetY:      offsetY,
		OffsetZ:      offsetZ,
		PlateType:    plate_type,
		Head:         head,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	r, err := llc.client.Aspirate(context.Background(), &pb.AspirateRequest{
		Volume:     volume,
		Overstroke: overstroke,
		Head:       head,
		Multi:      multi,
		Platetype:  platetype,
		What:       what,
		LLF:        llf,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	r, err := llc.client.Dispense(context.Background(), &pb.DispenseRequest{
		Volume:    volume,
		Blowout:   blowout,
		Head:      head,
		Multi:     multi,
		Platetype: platetype,
		What:      what,
		LLF:       llf,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	r, err := llc.client.LoadTips(context.Background(), &pb.LoadTipsRequest{
		Channels:  channels,
		Head:      head,
		Multi:     multi,
		Platetype: platetype,
		Position:  position,
		Well:      well,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	r, err := llc.client.UnloadTips(context.Background(), &pb.UnloadTipsRequest{
		Channels:  channels,
		Head:      head,
		Multi:     multi,
		Platetype: platetype,
		Position:  position,
		Well:      well,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	r, err := llc.client.SetPipetteSpeed(context.Background(), &pb.SetPipetteSpeedRequest{
		Head:    head,
		Channel: channel,
		Rate:    rate,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	r, err := llc.client.SetDriveSpeed(context.Background(), &pb.SetDriveSpeed{
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
		Head:      head,
		Volume:    volume,
		Platetype: platetype,
		Cycle:     cycle,
		Multi:     multi,
		What:      what,
		Blowout:   blowout,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) ResetPistons(head, channel int) driver.CommandStatus {
	r, err := llc.client.ResetPistons(context.Background(), &pb.ResetPistonsRequest{
		Head:    head,
		Channel: channel,
	})
	return llc.handleCommandReply(r, err)
}

func (llc *LowLevelClient) UpdateMetaData(props *LHProperties) driver.CommandStatus {
	if propsJSON, err := json.Marshal(props); err != nil {
		return driver.CommandStatus{Errorcode: driver.ERR, Msg: err.Error()}
	} else {
		r, err := llc.client.UpdateMetaData(context.Background(), &pb.UpdateMetaDataRequest{
			LHProperties_JSON: propsJSON,
		})
		return llc.handleCommandReply(r, err)
	}
}
