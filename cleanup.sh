#!/bin/bash
# cleanup of the ressources on absolute paths
res=(
    "/opt/hermes"
    "/etc/hermes"
    "/var/log/hermes"
    "/etc/systemd/system/hermes.service"
    "/etc/letsencrypt/renewal-hooks/post/hermes_start.sh"
    "/etc/letsencrypt/renewal-hooks/pre/hermes_stop.sh"
)
sudo rm -R ${res[*]}