#!/bin/bash
# Command to enable local testing
# Expect 
# ./run.sh prod
# Will load env in $ENVS_DIR/prod/studydb.env
env=$1
bin=./counter-service
set -a
source $ENVS_DIR/$env/studydb.env
$bin
