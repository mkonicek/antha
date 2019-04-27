#! /bin/sh

## This script is run by cloudbuild as part of the CI for antha itself.
set -e

COVERPKG=$(go list -f '{{ join .Deps "\n" }}' github.com/antha-lang/antha/... | sort | uniq | grep -i antha | tr '\n' ',' | sed -e 's/,$//')

exec go test -covermode=atomic -coverpkg="${COVERPKG}" -coverprofile=cover.out github.com/antha-lang/antha/...
