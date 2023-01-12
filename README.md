# hermes

[Github-Repository](https://github.com/frederik-jatzkowski/hermes)

A level-4-gateway for tls termination, load balancing and virtual hosting for tcp based services. It provides the following features:

- virtual hosting
- automated certificate management
- TLS-termination
- load balancing
- automated http-to-https redirects for better user experience

## Installation Guide

First, configure the a `.env` file with the following variables:

```
HERMES_ADMIN_HOST=<host>
HERMES_LOG_LEVEL=info
HERMES_EMAIL=<email>
HERMES_USER=<username>
HERMES_PASSWORD=<password>
```

Then, run hermes from the directory with the `.env`-file this using

```
docker run --rm --detach --network host \
	--volume hermes:/var/hermes/ \
	--volume letsencrypt:/etc/letsencrypt \
	--env-file .env \
	ghcr.io/frederik-jatzkowski/hermes:<tag>
```

## Configuration Guide

The hermes admin panel can be accessed at `https://${HERMES_ADMIN_HOST}:440`. You can log in using `HERMES_USER` and `HERMES_PASSWORD`. The herme configuration has the following nested structure:

### Gateway

These represent the TLS-Servers listening for incoming connections to forward. An configuration can theoretically have as many `Gateways` as there are adresses to listen on, but in most cases there will only be one or two, since multiplexing will mostly be done using the supplied [Server Name Indication](wikipedia.org/wiki/Server_Name_Indication) of incoming connections.

A `Gateway` has to be defined with a local address for the TLS-listener to listen on.

### Service

A `Gateway` can receive any number `Service`-nodes as children. These specify a service, that can be reached on this gateway and are identified by a servername. Hermes will automatically try to obtain a certificate for a service and renew it as necessary. Once a connection to hermes with a corresponding server name is established, hermes will perform a handshake and provide the necessary certificate.

A service in a gateway is defined by a server name, which has to be a valid [Server Name Indication](wikipedia.org/wiki/Server_Name_Indication). In addition, the machine must in fact be reachable using the servername. This is necessary to obtain and renew certificates with `certbot`.

### Server

Finally, we arrive at the definitions for the backend servers. Hermes will forward connections to them, that are intended for the service, they are defined in.

> Note, that hermes detects server outages and will proxy new connections only to servers that it has not marked as unavailable. In addition, hermes will regularly check on unavailable servers to see if they came back online. If you have for example a service with 10 backend servers, you can do maintainance on some of them without any effect for the user (as long as the remaining servers can handle the load).

Servers are defined by the remote address, the backend server can be reached at. Note, that any valid TCP/IP address can be used.

## Log Guide

The logs can be read using `docker logs`.
