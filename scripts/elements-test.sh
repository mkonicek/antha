#! /bin/sh

## This script is run by the gitlab ci for testing elements. It
## therefore makes various assumptions, and is not designed to be run
## locally. Proceed at your own risk.
set -e

branch=$1
commit_sha=$2

repo='
{
    "SchemaVersion": "2.0",
    "Repositories": {
        "repos.antha.com/antha-ninja/elements-westeros": {
            "Directory": "/elements",
            "Commit": "'$commit_sha'",
            "Branch": "'$branch'"
        }
    }
}
'

printf "$repo" > /tmp/repo.json
cat /tmp/repo.json
exec ./elements.test -test.timeout 1h -test.v /tmp/repo.json
