package logs

import "github.com/rs/zerolog"

const (
	Component     = "component"
	Server        = "server"
	Balancer      = "balancer"
	Service       = "service"
	Gateway       = "gateway"
	Port          = "port"
	ServerAddress = "serverAddress"
	ClientAddress = "clientAddress"
	HostName      = "hostName"
	Sni           = "sni"
)

func Trace() *zerolog.Event {
	return logger.Trace()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Fatal() *zerolog.Event {
	return logger.WithLevel(zerolog.FatalLevel)
}

func Panic() *zerolog.Event {
	return logger.WithLevel(zerolog.PanicLevel)
}
