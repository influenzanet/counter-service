#!/bin/bash
# Command to enable local testing
# Expect 
# ./run.sh prod
# Will load env in $ENVS_DIR/prod/studydb.env
env=$1
set -a
source $ENVS_DIR/$env/studydb.env
go run ./cmd/counter-service