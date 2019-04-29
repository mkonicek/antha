#! /bin/sh

## This script is run by cloudbuild as part of the CI for antha itself.
set -o nounset -o errexit -o pipefail -o noclobber
shopt -s failglob

## There are some packages that only contain test files. Go test gets
## upset if you try to include these packages in coverage, so we have
## to filter them out:
COVERPKG=$(go list -f '{{if (len .GoFiles) gt 0}}{{.ImportPath}}{{end}}' github.com/antha-lang/antha/... | tr '\n' ',' | sed -e 's/,$//')

go test -covermode=atomic -coverprofile=cover.out -coverpkg="${COVERPKG}" github.com/antha-lang/antha/...

COVERALLS_TOKEN_FILE=${HOME}/.coveralls_token
if [[ -f ${COVERALLS_TOKEN_FILE} && -f cover.out ]]; then
    COVERALLS_TOKEN=$(<${COVERALLS_TOKEN_FILE})
    goveralls -repotoken=${COVERALLS_TOKEN} -coverprofile=cover.out -package github.com/antha-lang/antha
fi
