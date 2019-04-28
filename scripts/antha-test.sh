#! /bin/sh

## This script is run by cloudbuild as part of the CI for antha itself.
set -e

COVERPKG=$(go list github.com/antha-lang/antha/... | tr '\n' ',' | sed -e 's/,$//')

exec go test -covermode=atomic -coverprofile=cover.out -coverpkg="${COVERPKG}" github.com/antha-lang/antha/...
