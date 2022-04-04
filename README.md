`v0.1.0-dev`

# hermes

A gateway for tls termination, load balancing and virtual hosting for tcp/tls based servers.

## Installation Guide

1. Make sure, that `certbot` is installed on your system:
        sudo apt-get update
        sudo apt-get install certbot

2. Download or copy the `hermes_linux_amd64.tar.gz` on the target system

3. Unpack the archive using `sudo tar -xvf hermes_linux_amd64.tar.gz -C /`

4. Configure your gateway using (for example) `sudo nano /etc/hermes/conf.xml`

5. Enable hermes using `sudo systemctl enable hermes.service`

6. Start hermes using `sudo systemctl start hermes.service`

7. Check the recent logs to see, if everything worked as expected using `sudo tail /var/log/hermes/logfile`

## Configuration Guide

### Restart after changes to the configuration

After changes have been made, restart hermes using ``

## Log Guide

The logs are located here: `/var/log/hermes/logfile`.