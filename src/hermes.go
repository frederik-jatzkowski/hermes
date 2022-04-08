package main

import (
	"encoding/xml"
	"os"
	"os/signal"
	"syscall"

	"fleo.software/infrastructure/hermes/certbot"
	"fleo.software/infrastructure/hermes/gateway"
	"fleo.software/infrastructure/hermes/http"
	"fleo.software/infrastructure/hermes/logs"
)

func main() {
	logs.BothPrint("hermes started", "0000")                  // log system start
	cfg := ParseConfig()                                      // read and parse cfg
	cfg.Init()                                                // initialize and validate cfg
	cfg.Run()                                                 // run cfg
	sigchan := make(chan os.Signal, 1)                        // make channel to await SIGTERM
	signal.Notify(sigchan, syscall.SIGTERM)                   //
	<-sigchan                                                 // await SIGTERM
	logs.ContinuousPrint("hermes gracefully stopped", "0010") // log program termination
	os.Exit(0)                                                // exit program successfully
}

type Config struct {
	Gateways []gateway.Gateway `xml:"Gateway"`
	Email    *string           `xml:"email,attr"`
	Ok       bool              `xml:"-"`
	Redirect *bool             `xml:"https-to-https,attr"`
}

func ParseConfig() *Config {
	cfg := &Config{
		Ok: true,
	} // create config
	file, err := os.ReadFile("/etc/hermes/config.xml") // read config file
	if err != nil {
		logs.LaunchPrint("could not find or read '/etc/hermes/config.xml' file", "1101")
		cfg.Ok = false // fatal
	} else {
		err = xml.Unmarshal(file, cfg) // try to parse
		if err != nil {
			logs.LaunchPrint("syntax error in '/etc/hermes/config.xml' file", "1201")
			cfg.Ok = false // fatal
		}
	}
	return cfg
}

func (cfg *Config) Init() {
	if cfg.Email == nil {
		logs.LaunchPrint("invalid config: missing email", "1301")
		cfg.Ok = false
	} else {
		certbot.SetEmail(*cfg.Email)
	}
	if !cfg.Ok {
		logs.BothPrint("invalid config: hermes could not start operating", "1001")
		return // only run, if valid
	}
	for i := range cfg.Gateways {
		cfg.Gateways[i].Init() // init all gateways
	}
	// start https redirect in case it is chosen
	if cfg.Redirect == nil || *cfg.Redirect {
		go http.ServeHTTPSRedirect()
	}
}

func (cfg *Config) Run() {
	if !cfg.Ok {
		return // only run, if valid
	}
	for i := range cfg.Gateways {
		go cfg.Gateways[i].Listen()
	}
}
