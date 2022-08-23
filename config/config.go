package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type HTTPSConfig struct {
	CertFile string
	KeyFile  string
}

type ServiceConfig struct {
	Target string
	Path   string
}

type DynamicServiceConfig struct {
	Path  string
	Query string
}

type Config struct {
	Service        []ServiceConfig
	DynamicService *DynamicServiceConfig
	HTTPS          *HTTPSConfig
}

var fileLock = &sync.RWMutex{}

func ReadConfig() *Config {
	fileLock.RLock()
	defer fileLock.RUnlock()
	file, _ := os.Open("./config.json")
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	decoder.Decode(&config)
	return &config
}

func WriteConfig(config string) {
	var out bytes.Buffer
	json.Indent(&out, []byte(config), "", "  ")
	fileLock.Lock()
	defer fileLock.Unlock()
	ioutil.WriteFile("./config.json", out.Bytes(), 0755)
}
