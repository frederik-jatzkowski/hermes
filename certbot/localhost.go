package certbot

import "crypto/tls"

var localhost tls.Certificate

func init() {
	localhost, _ = tlsLoadX509KeyPair("/etc/letsencrypt/live/localhost/fullchain.pem", "/etc/letsencrypt/live/localhost/privkey.pem")
}
