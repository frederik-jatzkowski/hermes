`v0.1.0-dev`

# hermes

[Github-Repository](https://github.com/frederik-jatzkowski/hermes)

A gateway for tls termination, load balancing and virtual hosting for tcp or tls based services. It provides the following features:

- virtual hosting of tcp or tls based services that will be visible over tls
- automated management of certificates for the required server names
- TLS-termination
- load balancing of connections to multiple backend servers
- encrypted and unencrypted connections to backend servers
- automated http-to-https redirects for better user experience

## Installation Guide

1. Download or copy the `hermes_linux_amd64.tar.gz` on the target system.

2. Unpack the archive using `sudo tar -xvf hermes_linux_amd64.tar.gz -C /` so all files are installed in the correct location.

3. Configure your gateway using (for example) `sudo nano /etc/hermes/config.xml`.

4. Enable hermes using `sudo systemctl enable hermes.service`.

5. Start hermes using `sudo systemctl start hermes.service`.

6. Check the recent logs to see, if any errors occurred during the launch of hermes `sudo tail /var/log/hermes/launch.log`

7. See `sudo cat /var/log/hermes/operation.log` for errors and info about the running system

## Configuration Guide

Hermes uses XML as the configuration format since it resembles the nested architecture of the system. It read the configuration from `/etc/hermes/config.xml`. In the follwing sections we will build a configuration file together.

> Note: After changes to the config file have been made, restart hermes using `sudo systemctl restart hermes` to apply changes.

### Config

The root node of the configuration has to be a `Config`-node:
```
<?xml version="1.0" encoding="UTF-8"?>

<Config email="admin@my.org">
    ...
</Config>
```

The attribute `email` is necessary for obtaining and renewing the necessary certificates. It will be handed over to `certbot` to deliver urgent information about the certificates.

Optionally, a boolean attribute named `http-to-https` can be provided. Its default value is "true". This tells hermes to spin up a simple http-Server on the standard http-port (usually `:80`).
This server will respond with a status code of `301 Moved Permanently` and redirect traffic to `https`. This behaviour can be avoided by settings the attribute `http-to-https="false"` on the `Config`-node.


### Gateway

The children of the `Config`-node are `Gateway`-nodes. These represent the TLS-Servers listening for incoming connections to forward. A `Config` can theoretically have as many `Gateways` as there are adresses to listen on, but in most cases there will only be one or two, since multiplexing will mostly be done using the supplied [Server Name Indication](wikipedia.org/wiki/Server_Name_Indication) of incoming connections.

In the follwing example we define a `Gateway` that listens on port `443`. This is the default port for `https`-traffic thus the `Gateway` will receive all incoming connections, that try to access a `https`-server at this IP.

```
<?xml version="1.0" encoding="UTF-8"?>

<Config email="admin@my.org">
    <Gateway laddress="0.0.0.0:443">
		...
	</Gateway>
</Config>
```

As you can see, a `Gateway` has to be defined with the attribute `laddress`, which has to be a local address for the TLS-listener to listen on.

### Service

A `Gateway` can receive any number `Service`-nodes as children. These specify a service, that can be reached on this gateway and are identified by a servername. Hermes will automatically try to obtain a certificate for a service and renew it as necessary. Once a connection to hermes with a corresponding server name is established, hermes will perform a handshake and provide the necessary certificate.

Let's add a service to our example:

```
<?xml version="1.0" encoding="UTF-8"?>

<Config email="admin@my.org">
    <Gateway laddress="0.0.0.0:443">
		<Service servername="localhost">
		    ...
		</Service>
	</Gateway>
</Config>
```

The attribute `servername` has to be a valid [Server Name Indication](wikipedia.org/wiki/Server_Name_Indication). In addition, the machine must in fact be reachable using the servername. This is necessary to obtain and renew certificates with `certbot`.


### LoadBalancer

For simplicity and thus reliablity, every `Service` node must have exactly one `LoadBalancer` child node. Even if there is only one backend server for the service, a `LoadBalancer` has to be specified.

Let's add one:

```
<?xml version="1.0" encoding="UTF-8"?>

<Config email="admin@my.org">
    <Gateway laddress="0.0.0.0:443">
		<Service servername="localhost">
		    <LoadBalancer>
				...
			</LoadBalancer>
		</Service>
	</Gateway>
</Config>
```

Optionally, you can specify the `algorithm`-attribute and select the algorithm used for balancing. In this version of hermes, only one algorithm, `"RoundRobin"` is available. This is also the default value.


### Server

Finally, we arrive at the definitions for the backend servers. These are defined inside a `LoadBalancer` and hermes will forward connections to them, that are intended for the service, they are defined in.

> Note, that hermes detects server outages and will proxy new connections only to servers that it has not marked as unavailable. In addition, hermes will regularly check on unavailable servers to see if they came back online. If you have for example a service with 10 backend servers, you can do maintainance on some of them without any effect for the user (as long as the remaining servers can handle the load).

In the example configuration, this might look like this:

```
<?xml version="1.0" encoding="UTF-8"?>

<Config email="admin@my.org">
    <Gateway laddress="0.0.0.0:443">
		<Service servername="localhost">
		    <LoadBalancer>
				<Server raddress="secure:port" tls="false"/>
				<Server raddress="insecure:port" tls="true"/>
			</LoadBalancer>
		</Service>
	</Gateway>
</Config>
```

The attribute `raddress` specifies the address, the backend server can be reached at. Note, that any valid TCP/IP address can be used.

The optionaly second attribute is the boolean `tls`-attribute. This specifies, wether the connection between hermes and the server should be encrypted with TLS or send as unencrypted TCP. Default is `"false"`.

### Conclusion

The configuration that we built would be fully functional, if the correct addresses, server name and email would be filled in. For hermes to work nominally, its machine has to be reachable on the specified server names and a sufficient amount of backend servers have to be up and running for each service.


## Log Guide

The logs for the running system are located here: `/var/log/hermes/operation.log`.

The logs for launches of the system are located here: `/var/log/hermes/launch.log`.

### Log Event Codes

In this section, the event codes associated with logged events are listed and explained. Event codes have the format `<module>-xxx`, where `x` is a number between `0` and `9` and `<module>` specifies the module that the event happened in.

In general, success-events will have an associated event code that ends with a `0`. Codes with another ending will generally mean, that something did not work as expected and the functionality of hermes has been limited by the event.

## Peformance and Reliability

During tests, a hermes instance on a tiny virtual server with one vCPU reliably handled around 180 HTTPS-requests with a 4kb response every second and distributed them on 10 different servers on the same vm. On a machine with multiple cores and seperated backend servers, this value will likely be higher since hermes then can make better use of golangs powerful goroutines, although a proper test series has not been performed. During a stresstest of 20,000 requests within 2 minutes, not a single request failed. The theoretical limit of concurrent connections (not requests) is the number of available ports on the system and thereby lies around 60,000.