`v0.1.0-dev`

# hermes

A gateway for tls termination, load balancing and virtual hosting for tcp/tls based servers.


## Installation Guide

1. Download or copy the `hermes_linux_amd64.tar.gz` on the target system.

2. Unpack the archive using `sudo tar -xvf hermes_linux_amd64.tar.gz -C /` so all files are installed in the correct location.

3. Configure your gateway using (for example) `sudo nano /etc/hermes/conf.xml`.

4. Enable hermes using `sudo systemctl enable hermes.service`.

5. Start hermes using `sudo systemctl start hermes.service`.

6. Check the recent logs to see, if any errors occurred during the launch of hermes `sudo tail /var/log/hermes/launch.log`

7. See `sudo cat /var/log/hermes/operation.log` for errors and info about the running system

## Configuration Guide

Note: After changes to the config file have been made, restart hermes using `sudo systemctl restart hermes` to apply changes.

### Config

### Gateway

### Service

### LoadBalancer

### Server

## Log Guide

The logs for the running system are located here: `/var/log/hermes/operation.log`.

The logs for launches of the system are located here: `/var/log/hermes/launch.log`.