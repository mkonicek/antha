generate:
	go generate
	find . -name '*.pb.go' | xargs perl -p -i \
	  -e 's|github.com/golang/protobuf/proto|github.com/antha-lang/antha/bvendor/github.com/golang/protobuf/proto|g;' \
	  -e 's|golang.org/x/net/context|github.com/antha-lang/antha/bvendor/golang.org/x/net/context|g;' \
	  -e 's|google.golang.org/grpc|github.com/antha-lang/antha/bvendor/google.golang.org/grpc|g;'
