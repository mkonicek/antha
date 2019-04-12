#! /bin/bash

set -ex
repos="github.com/Synthace/instruction-plugins github.com/Synthace/antha-runner"
for r in $repos; do
    go get $r
    ( cd /go/src/$r && git fetch --all && git checkout feature/future_sanity )
done
