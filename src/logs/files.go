package logs

import (
	"log"
	"os"
)

var logfile *os.File = prepareFile("/var/log/hermes/operation.log")
var launchfile *os.File = prepareFile("/var/log/hermes/launch.log")

func prepareFile(name string) *os.File {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 700)
	if err != nil {
		log.Fatal(err, "file could not be accessed", name)
	}
	return f
}
