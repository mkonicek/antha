SHELL=/bin/bash
ASL=antha/AnthaStandardLibrary/Packages

all: gen_comp compile

gen_comp:
	go run cmd/antha/antha.go -outdir=antha/component/lib antha/component/an
	gofmt -w -s antha/component/lib

test:
	go test -v `go list ./... | grep -v vendor | grep -v bvendor`

gen_pb:
	go generate github.com/antha-lang/antha/driver
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|proto "([^"]*)"|proto "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|context "([^"]*)"|context "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|grpc "([^"]*)"|grpc "github.com/antha-lang/antha/bvendor/\1"|'

fmt_json:
	for i in `find antha/examples -name '*.json' -o -name '*.yml'`; do \
	  python -mjson.tool "$$i" > "$$i.bak" && mv "$$i.bak" "$$i"; \
	done

compile:
	go install github.com/antha-lang/antha/cmd/...

test_workflows:
	go test github.com/antha-lang/antha/antha/component/lib

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

.PHONY: all gen_comp fmt_json test test_workflows compile assets ALWAYS
