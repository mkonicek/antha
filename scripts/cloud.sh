#!/bin/bash -eu

## Expects the environment variable $WF_JSON to be set to a workflow value.
## Executes $* but with the given workflow inserted as a final file argument.
exec "$@" <(echo "$WF_JSON")
