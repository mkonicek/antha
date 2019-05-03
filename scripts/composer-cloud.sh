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
trap "{ > $DATA_DIR/outputReady; exit 0; }" EXIT

echo input
find ${DATA_DIR}/input -type f
$MIGRATE -from=${DATA_DIR}/input/workflow/request.pb -outdir=${DATA_DIR}/scratch -gilson-device=gillian -format=protobuf - <<<$WF_JSON

echo scratch
find ${DATA_DIR}/scratch -type f
cp -a ${DATA_DIR}/scratch/* ${DATA_DIR}/input

echo input
find ${DATA_DIR}/input -type f
$COMPOSER -indir=${DATA_DIR}/scratch -outdir=${DATA_DIR}/output -linkedDrivers

echo output
find ${DATA_DIR}/output -type f
