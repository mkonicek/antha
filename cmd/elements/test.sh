#! /bin/sh

## This script is run by the gitlab ci for testing elements. It
## therefore makes various assumptions, and is not designed to be run
## locally. Proceed at your own risk.
set -e

commit_sha=$1

repo='
{
    "SchemaVersion": "2.0",
    "Repositories": {
        "repos.antha.com/antha-ninja/elements-westeros": {
            "Directory": "/elements",
            "Commit": "'$commit_sha'"
        }
    }
}
'

printf "$repo" > /tmp/repo.json
cat /tmp/repo.json
exec go test github.com/antha-lang/antha/cmd/elements -v -args /tmp/repo.json
