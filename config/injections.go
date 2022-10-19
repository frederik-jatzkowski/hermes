package config

import "os"

var osReadFile func(name string) ([]byte, error) = os.ReadFile
