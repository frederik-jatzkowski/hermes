#!/bin/bash

# env GOOS=linux GOARCH=arm64 go build -o /run/hermes/hermes ./hermes.go

# compile arm64
cd src
env GOOS=linux GOARCH=amd64 go build -o ../dist/run/hermes/hermes ./hermes.go
cd ../dist
# distribute ressources to absolute paths
sudo cp -R ./etc /
sudo cp -R ./run /
sudo cp -R ./var /
# archive all ressources
tar -cvzf ./hermes_linux_amd64.tar.gz /run/hermes /var/log/hermes /etc/systemd/system/hermes.service /etc/hermes
# cleanup of the ressources on absolute paths
sudo rm -R /etc/hermes
sudo rm -R /run/hermes
sudo rm -R /var/log/hermes
sudo rm -R /etc/systemd/system/hermes.service
cd ..