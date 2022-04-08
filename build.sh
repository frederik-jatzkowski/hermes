#!/bin/bash

# env GOOS=linux GOARCH=arm64 go build -o /run/hermes/hermes ./hermes.go

# compile arm64
cd src
if env GOOS=linux GOARCH=amd64 go build -o ../dist/opt/hermes/hermes ./hermes.go ; then
    cd ../dist
    # distribute ressources to absolute paths
    sudo cp -R ./etc /
    sudo cp -R ./opt /
    sudo cp -R ./var /
    # array of ressources
    res=(
        "/opt/hermes/hermes"
        "/etc/hermes/config.xml"
        "/var/log/hermes/launch.log"
        "/var/log/hermes/operation.log"
        "/etc/systemd/system/hermes.service"
        "/etc/letsencrypt/renewal-hooks/post/hermes_start.sh"
        "/etc/letsencrypt/renewal-hooks/pre/hermes_stop.sh"
    )
    # archive all ressources
    sudo tar cvzfP ./hermes_linux_amd64.tar.gz ${res[*]}
    # cleanup of the ressources on absolute paths
    sudo rm -R ${res[*]}
fi
cd ..