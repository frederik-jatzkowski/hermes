#!/bin/bash

# env GOOS=linux GOARCH=arm64 go build -o /run/hermes/hermes ./hermes.go

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
arch=amd64 # build arm64
cd src
if env GOOS=linux GOARCH=$arch go build -o ../dist/opt/hermes/hermes ./hermes.go ; then
    cd ../dist
    sudo cp -R ./etc ./opt ./var / # distribute ressources to absolute paths
    sudo tar cvzfP ./hermes_linux_$arch.tar.gz ${res[*]} # archive all ressources
    sudo rm -R ${res[*]} # cleanup of the ressources on absolute paths
fi
cd ..
arch=arm64 # build arm64
cd src
if env GOOS=linux GOARCH=$arch go build -o ../dist/opt/hermes/hermes ./hermes.go ; then
    cd ../dist
    sudo cp -R ./etc ./opt ./var / # distribute ressources to absolute paths
    sudo tar cvzfP ./hermes_linux_$arch.tar.gz ${res[*]} # archive all ressources
    sudo rm -R ${res[*]} # cleanup of the ressources on absolute paths
fi
cd ..