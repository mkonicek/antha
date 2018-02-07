package driver

import "github.com/golang/protobuf/proto"

// A Call is a generic call to a device
type Call struct {
	Method string
	Args   proto.Message
	Reply  proto.Message
}
