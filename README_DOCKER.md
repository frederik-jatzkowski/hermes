```
mkdir /etc/hermes
touch /etc/hermes/config.xml
mkdir /var/log/hermes
touch /etc/hermes/hermes.log
docker run --mount type=bind,src=/var/log/hermes,dst=/var/log/hermes,readonly
```

```
HERMES_IP=<ip>
HERMES_EMAIL=<email>

docker run --env-file="./.env" --network="host" hermes
docker run --network="host" hermes -e \${HERMES_EMAIL} -a \${HERMES_ADMIN_HOST} -l \${HERMES_LOG_LEVEL} -u \${HERMES_USER} -p \${HERMES_PASSWORD}
docker run --network="host" hermes -e ${HERMES_EMAIL} -a ${HERMES_ADMIN_HOST} -l ${HERMES_LOG_LEVEL} -u ${HERMES_USER} -p ${HERMES_PASSWORD}
```
