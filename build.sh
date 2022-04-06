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
        "/opt/hermes"
        "/etc/hermes"
        "/var/log/hermes"
        "/etc/systemd/system/hermes.service"
        "/etc/letsencrypt/renewal-hooks/post/hermes_start.sh"
        "/etc/letsencrypt/renewal-hooks/pre/hermes_stop.sh"
    )
    # archive all ressources
    sudo tar cvzfP ./hermes_linux_amd64.tar.gz ${res[*]}
    # cleanup of the ressources on absolute paths
    sudo rm -R ${res[*]}
    # sudo rm -R /etc/hermes
    # sudo rm -R /etc/letsencrypt/renewal-hooks/pre/hermes_stop.sh
    # sudo rm -R /etc/letsencrypt/renewal-hooks/post/hermes_start.sh
    # sudo rm -R /opt/hermes
    # sudo rm -R /var/log/hermes
    # sudo rm -R /etc/systemd/system/hermes.service
fi
cd ..