#!/bin/bash

# tell bash to exit on error
set -e

# disable bash history
set +o history

# setup of env-variables
export $(grep -v '^#' .env | xargs)

# build docker container
docker build . -t hermes:latest