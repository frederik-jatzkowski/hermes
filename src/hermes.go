package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"fleo.software/infrastructure/hermes/gateway"
	"fleo.software/infrastructure/hermes/logging/startup"
)

func main() {
	log.Println("program start")
	// read config file
	//file, err := os.ReadFile("./hermes.config.xml")
	file, err := os.ReadFile("./hermes.config.xml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// unmarshal config file into cfg
	cfg := &Config{}
	err = xml.Unmarshal(file, cfg)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalln("invalid config, program terminated")
	}
	// initialize and validate cfg
	collector := cfg.Init()
	if !collector.IsEmpty() {
		collector.Print()
		log.Fatalln("invalid config, program terminated")
	}
	// success message
	go cfg.Run()
	fmt.Println("configuration valid. starting gateway...")
	// await SIGTERM
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM)
	<-sigchan
	log.Println("graceful termination")
}

type Config struct {
	Gateways []gateway.Gateway `xml:"Gateway"`
}

func (cfg *Config) Init() *startup.ErrorCollector {
	collector := startup.NewErrorCollector()
	for i := 0; i < len(cfg.Gateways); i++ {
		cfg.Gateways[i].Init(collector)
	}
	return collector
}

func (cfg *Config) Run() {
	var wg sync.WaitGroup
	for i := 0; i < len(cfg.Gateways); i++ {
		wg.Add(1)
		gw := cfg.Gateways[i]
		go func() {
			defer wg.Done()
			gw.Listen()
		}()
	}
	wg.Wait()
}
