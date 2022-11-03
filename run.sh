#!/bin/bash

# tell bash to exit on error
set -e

# disable bash history
set +o history

# setup of env-variables
export $(grep -v '^#' .env | xargs)

# run container
# docker run --network="host" --rm -v hermes:/var/hermes/ hermes -e ${HERMES_EMAIL} -l ${HERMES_LOG_LEVEL} -u ${HERMES_USER} -p ${HERMES_PASSWORD} #-a ${HERMES_ADMIN_HOST} 
docker run --network host --rm --volume hermes:/var/hermes/ --env-file .env hermes