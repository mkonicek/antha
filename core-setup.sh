#! /bin/bash

set -ex
repos="github.com/Synthace/instruction-plugins github.com/Synthace/antha-runner"
for r in $repos; do
    mkdir -p /go/src/$r
    ( cd /go/src/$r && git clone https://$r . && git checkout feature/future_sanity )
done

for r in $repos; do
    go get ${r}...
done
