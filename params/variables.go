package params

import (
	"time"
)

// cli params
var Version string = "dev"
var EmailAdress string
var User string
var Password string
var HealthCheckInterval time.Duration = time.Second
var CertCheckInterval time.Duration = time.Hour * 24
var AdminHost string
var LogLevel string
