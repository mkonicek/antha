SHELL=/bin/bash
ASL=antha/AnthaStandardLibrary/Packages

all:

test:
	go test `go list ./... | grep -v vendor | grep -v bvendor`

gen_pb:
	go generate github.com/antha-lang/antha/driver
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|proto "([^"]*)"|proto "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|context "([^"]*)"|context "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|grpc "([^"]*)"|grpc "github.com/antha-lang/antha/bvendor/\1"|'

assets: $(ASL)/asset/asset.go

$(ASL)/asset/asset.go: $(GOPATH)/bin/go-bindata-assetfs $(ASL)/asset_files/rebase/type2.txt
	cd $(ASL)/asset_files && $(GOPATH)/bin/go-bindata-assetfs -pkg=asset ./...
	mv $(ASL)/asset_files/bindata_assetfs.go $@
	gofmt -s -w $@

$(ASL)/asset_files/rebase/type2.txt: ALWAYS
	mkdir -p `dirname $@`
	curl -o $@ ftp://ftp.neb.com/pub/rebase/type2.txt

$(GOPATH)/bin/2goarray:
	go get -u github.com/cratonica/2goarray

$(GOPATH)/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

$(GOPATH)/bin/go-bindata-assetfs: $(GOPATH)/bin/go-bindata
	go get -u -f github.com/elazarl/go-bindata-assetfs/...
	touch $@

.PHONY: all test assets ALWAYS
