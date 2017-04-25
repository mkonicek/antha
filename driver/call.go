package driver

import "github.com/golang/protobuf/proto"

type Call struct {
	Method string
	Args   proto.Message
	Reply  proto.Message
}
