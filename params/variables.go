package params

import (
	"os"
	"time"
)

// cli params or env variables
var Version string = "dev"
var EmailAdress string = os.Getenv("HERMES_EMAIL")
var User string = os.Getenv("HERMES_USER")
var Password string = os.Getenv("HERMES_PASSWORD")
var AdminHost string = os.Getenv("HERMES_ADMIN_HOST")
var LogLevel string = os.Getenv("HERMES_LOG_LEVEL")
var HealthCheckInterval time.Duration = time.Second
var CertCheckInterval time.Duration = time.Hour * 24
