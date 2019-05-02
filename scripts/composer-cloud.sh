#!/bin/bash -eu
# Ambassador-compatible wrapper script for composer. Expects the environment
# variable $WF_JSON to be set to a workflow value.  Other flags could also be
# set, but not -indir and -outdir.
COMPOSER=${COMPOSER:-composer}
MIGRATE=${MIGRATE:-migrate}

# Well-known directory name set by the workload service.
# (See https://github.com/Synthace/microservice/tree/master/cmd/workload#compatible-containers )
DATA_DIR=${DATA_DIR:-/data}

# (See https://github.com/Synthace/microservice/tree/master/cmd/ambassador )
{ while ! [[ -r $DATA_DIR/inputReady ]]; do sleep 0.1; done; }
< $DATA_DIR/inputReady
trap "{ > $DATA_DIR/outputReady; }" EXIT

$MIGRATE -from=${DATA_DIR}/input/request.pb -outdir=${DATA_DIR}/scratch -gilson-device=gillian -format=protobuf - <<<$WF_JSON

$COMPOSER -indir=${DATA_DIR}/scratch -outdir=${DATA_DIR}/output -linkedDrivers
