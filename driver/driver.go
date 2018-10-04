//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_runner_v1/runner.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_driver_v1/driver.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_shakerincubator_v1/shakerincubator.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_human_v1/human.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_platereader_v1/platereader.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_framework_v1/framework.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I${GOPATH}/src ${GOPATH}/src/github.com/antha-lang/antha/driver/antha_quantstudio_v1/quantstudio.proto --go_out=plugins=grpc:${GOPATH}/src
//go:generate protoc -I. lh/lh.proto --go_out=plugins=grpc:pb

package driver
