SHELL=/bin/bash
ASL=antha/AnthaStandardLibrary/Packages
PKG=$(shell go list .)

# Check if we support verbose
XARGS_HAS_VERBOSE=$(shell if echo | xargs --verbose 2>/dev/null; then echo -n '--verbose'; fi)

all:

build:
	go install ./cmd/antha

# Build first to speed up tests
test: build
	go test ./...

lint: test
	golangci-lint run --deadline=5m -E gosec -E gofmt ./...

docker-build: .build/antha-build-image .build/antha-build-withdeps-image

.build/antha-build-image:
	mkdir -p .build
	docker build -t antha-build .
	touch $@

.build/antha-build-withdeps-image: .build/antha-build-image
	mkdir -p .build
	docker run --rm -v `pwd`:/go/src/$(PKG) -w /go/src/$(PKG) \
	  antha-build make .build/imports
	docker build -f Dockerfile.withdeps -t antha-build-withdeps .
	touch $@

.build/imports:
	(go list -f '{{join .Imports "\n"}}' ./... && go list -f '{{join .TestImports "\n"}}' ./...) \
	  | sort | uniq | grep -v $(PKG) > $@

docker-lint: .build/antha-build-withdeps-image
	docker run --rm -v `pwd`:/go/src/$(PKG) -w /go/src/$(PKG) \
	  antha-build-withdeps make lint

gen_pb:
	go generate $(PKG)/driver

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

.PHONY: all test lint docker_lint get_deps assets ALWAYS
