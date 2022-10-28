package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

var access = sync.Mutex{}

const configsDirectory = "/var/hermes/configs/"

func ConfigHistory() ([]Config, error) {
	var (
		dir     []fs.DirEntry
		data    []byte
		configs []Config
		err     error
	)

	// lock configs
	access.Lock()
	defer access.Unlock()

	// read all configs in config directory
	dir, err = os.ReadDir(configsDirectory)
	if err != nil {
		return configs, fmt.Errorf("could not read configs directory: %s", err)
	}

	// read all files
	for _, entry := range dir {
		// check for unix timestamp as filename
		unix, err := strconv.ParseInt(entry.Name(), 10, 64)
		if err != nil {
			continue
		}

		// read config file
		data, err = os.ReadFile(configsDirectory + entry.Name())
		if err != nil {
			return configs, fmt.Errorf("could not read file: %s", err)
		}

		// parse and validate config
		config, err := NewConfig(data)
		if err != nil {
			return configs, fmt.Errorf("could not build config: %s", err)
		}

		// use unix timestamp from filename
		config.Unix = unix

		// approve config
		configs = append(configs, config)
	}

	// sort configs by their unix timestamp
	sort.Slice(configs, func(i int, j int) bool {
		return configs[i].Unix < configs[j].Unix
	})

	return configs, err
}

func AppendConfig(config Config) error {
	// lock access
	access.Lock()
	// wait for more than 1 second after access is given back to the configs directory
	defer time.Sleep(time.Millisecond * 1500)
	defer access.Unlock()

	// use a timestamp 1 second in the future
	unix := time.Now().Unix() + 1
	filename := strconv.FormatInt(unix, 10)
	file, err := os.OpenFile(configsDirectory+filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create file: %s", err)
	}

	// marshal config
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config: %s", err)
	}

	// write config to disk
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("could not write config to file: %s", err)
	}

	return err
}
