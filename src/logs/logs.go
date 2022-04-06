package logs

import (
	"log"
)

var continuousLogger log.Logger = *log.New(logfile, "", log.Flags())
var launchLogger log.Logger = *log.New(launchfile, "", log.Flags())

func ContinuousPrint(msg interface{}, code string) {
	continuousLogger.Print(msg, " (event code "+code+")")
}

func LaunchPrint(msg interface{}, code string) {
	launchLogger.Print(msg, " (event code "+code+")")
}

func BothPrint(msg interface{}, code string) {
	ContinuousPrint(msg, code)
	LaunchPrint(msg, code)
}
