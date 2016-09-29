package driver

import "github.com/antha-lang/antha/bvendor/github.com/golang/protobuf/proto"

type Call struct {
	Method string
	Args   proto.Message
	Reply  proto.Message
}
