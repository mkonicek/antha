FROM golang:1.6-wheezy
RUN apt-get update && apt-get install -y libglpk-dev libsqlite3-dev libopencv-dev
ADD . /go/src/github.com/antha-lang/antha
RUN go get golang.org/x/net/context
RUN go install github.com/antha-lang/antha/cmd/...
