#! /bin/sh

## This script is run by cloudbuild as part of the CI for antha itself.
set -e

## There are some packages that only contain test files. Go test gets
## upset if you try to include these packages in coverage, so we have
## to filter them out:
COVERPKG=$(go list -f '{{if (len .GoFiles) gt 0}}{{.ImportPath}}{{end}}' github.com/antha-lang/antha/... | tr '\n' ',' | sed -e 's/,$//')

exec go test -covermode=atomic -coverprofile=cover.out -coverpkg="${COVERPKG}" github.com/antha-lang/antha/...
