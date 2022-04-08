#!/bin/bash


systemctl stop hermes
systemctl disable hermes

res=(
    "/opt/hermes/hermes"
    "/etc/hermes/config.xml"
    "/var/log/hermes/launch.log"
    "/var/log/hermes/operation.log"
    "/etc/systemd/system/hermes.service"
    "/etc/letsencrypt/renewal-hooks/post/hermes_start.sh"
    "/etc/letsencrypt/renewal-hooks/pre/hermes_stop.sh"
)
rm -R ${res[*]}

tar -xvf hermes_linux_amd64.tar.gz -C /

cat ./my-config.xml 1> /etc/hermes/config.xml

systemctl enable hermes && systemctl start hermes