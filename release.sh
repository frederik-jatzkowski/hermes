#!/bin/bash

# tell bash to exit on error
set -e

# disable bash history
set +o history

# setup of env-variables
export $(grep -v '^#' .env | xargs)

# build docker container
docker build . -t hermes:latest

# run container
docker run --network="host" hermes -e ${HERMES_EMAIL} -a ${HERMES_ADMIN_HOST} -l ${HERMES_LOG_LEVEL} -u ${HERMES_USER} -p ${HERMES_PASSWORD}