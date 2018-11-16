package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	pb "github.com/antha-lang/antha/driver/liquidhandling/pb"
	driver "github.com/antha-lang/antha/microArch/driver"
	liquidhandling "github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type HighLevelClient struct {
	client pb.HighLevelLiquidhandlingDriverClient
}

// NewHighLevelClient create a client for connecting with a remote high level
// server
func NewHighLevelClient(address string) (*HighLevelClient, error) {
	if conn, err := grpc.Dial(address, grpc.WithInsecure()); err != nil {
		return nil, errors.WithMessage(err, "Cannot initialize driver")
	} else {
		return NewHighLevelClientFromConn(conn), nil
	}
}

// NewHighLevelClientFromConn create a client for connecting with a remote high level server from a grpc Conn object
func NewHighLevelClientFromConn(conn *grpc.ClientConn) *HighLevelClient {
	return &HighLevelClient{
		client: pb.NewHighLevelLiquidhandlingDriverClient(conn),
	}
}

func (hlc *HighLevelClient) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	if obj, ok := plate.(wtype.LHObject); !ok {
		return driver.CommandStatus{
			Errorcode: driver.ERR,
			Msg:       fmt.Sprintf("unable to serialize object of type %T", plate),
		}
	} else if plateJSON, err := wtype.MarshalDeckObject(obj); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else if r, err := hlc.client.AddPlateTo(context.Background(), &pb.AddPlateToRequest{
		Position:   position,
		Plate_JSON: string(plateJSON),
		Name:       name,
	}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) RemoveAllPlates() driver.CommandStatus {
	if r, err := hlc.client.RemoveAllPlates(context.Background(), &pb.RemoveAllPlatesRequest{}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) RemovePlateAt(position string) driver.CommandStatus {
	if r, err := hlc.client.RemovePlateAt(context.Background(), &pb.RemovePlateAtRequest{Position: position}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) Initialize() driver.CommandStatus {
	if r, err := hlc.client.Initialize(context.Background(), &pb.InitializeRequest{}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) Finalize() driver.CommandStatus {
	if r, err := hlc.client.Finalize(context.Background(), &pb.FinalizeRequest{}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	if r, err := hlc.client.Message(context.Background(), &pb.MessageRequest{
		Level:      int32(level),
		Title:      title,
		Text:       text,
		ShowCancel: showcancel,
	}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}

func (hlc *HighLevelClient) GetOutputFile() ([]byte, driver.CommandStatus) {
	if r, err := hlc.client.GetOutputFile(context.Background(), &pb.GetOutputFileRequest{}); err != nil {
		return nil, driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return r.OutputFile, driver.CommandStatus{OK: r.Status.OK, Errorcode: int(r.Status.Errorcode), Msg: r.Status.Msg}
	}
}

func (hlc *HighLevelClient) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	if r, err := hlc.client.GetCapabilities(context.Background(), &pb.GetCapabilitiesRequest{}); err != nil {
		return liquidhandling.LHProperties{}, driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		var ret liquidhandling.LHProperties
		if err := json.Unmarshal([]byte(r.LHProperties_JSON), &ret); err != nil {
			return ret, driver.CommandStatus{Errorcode: driver.ERR, Msg: err.Error()}
		}
		return ret, driver.CommandStatus{OK: r.Status.OK, Errorcode: int(r.Status.Errorcode), Msg: r.Status.Msg}
	}
}

func (hlc *HighLevelClient) Transfer(what, platefrom, wellfrom, plateto, wellto []string, volume []float64) driver.CommandStatus {
	if r, err := hlc.client.Transfer(context.Background(), &pb.TransferRequest{
		What:      what,
		Platefrom: platefrom,
		Wellfrom:  wellfrom,
		Plateto:   plateto,
		Wellto:    wellto,
		Volume:    volume,
	}); err != nil {
		return driver.CommandStatus{
			Msg:       err.Error(),
			Errorcode: driver.ERR,
		}
	} else {
		return driver.CommandStatus{OK: r.OK, Errorcode: int(r.Errorcode), Msg: r.Msg}
	}
}
