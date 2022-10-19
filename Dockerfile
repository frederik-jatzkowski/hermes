# multi stage build
FROM golang:1.18-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build --ldflags="-X github.com/frederik-jatzkowski/hermes/params.Version=${HERMES_VERSION}" -o /out/

FROM alpine

RUN apk add --no-cache certbot
#RUN apt-get update && apt-get install certbot -y

COPY --from=build /out/hermes /opt/hermes/hermes
COPY ./configs.json /var/hermes/configs.json
COPY ./hermes.log /var/hermes/hermes.log
COPY ./localhost/* /etc/letsencrypt/live/localhost/

VOLUME [ "/var/hermes" ]

#ENTRYPOINT /opt/hermes/hermes --version
ENTRYPOINT [ "/opt/hermes/hermes" ]
#ENTRYPOINT ["/opt/hermes/hermes", "-e", "\${HERMES_EMAIL}", "-a", "\${HERMES_ADMIN_HOST}", "-l", "\${HERMES_LOG_LEVEL}", "-u", "\${HERMES_USER}", "-p", "\${HERMES_PASSWORD}"]