#! /bin/sh

## This script is run by the gitlab ci for uploading element sets. It
## therefore makes various assumptions, and is not designed to be run
## locally. Proceed at your own risk.
set -e

endpoint=$1
branch=$2
commit_sha=$3
auth_token=$4

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

elements describe -format=protobuf /tmp/repo.json > /tmp/elements.pb
elements defaults /tmp/repo.json > /tmp/metadata.json
upload-metadata -element-proto /tmp/elements.pb -metadata-json /tmp/metadata.json \
    -endpoint="$endpoint" -element-set="$branch" \
    -token="$auth_token"
