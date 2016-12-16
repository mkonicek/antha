package api

type GrpcMessage struct {
	Name string `json:"name"` // Fully qualified message name
	Data []byte `json:"data"` // Protobuf data
}

type GrpcCall struct {
	Method string       `json:"method"` // Fully qualified rpc call name
	Args   *GrpcMessage `json:"args"`
	Reply  *GrpcMessage `json:"reply"`
}
