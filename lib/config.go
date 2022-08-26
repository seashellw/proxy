package lib

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

type FileServiceConfig struct {
	Path string
	Dir  string
}

type Config struct {
	Password       string
	Service        []ServiceConfig
	FileService    []FileServiceConfig
	DynamicService *DynamicServiceConfig
	HTTPS          *HTTPSConfig
}

var fileLock = &sync.RWMutex{}

var ConfigFilePath = "./config.json"

func (config *Config) ReadConfig() {
	fileLock.RLock()
	defer fileLock.RUnlock()
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&config)
}

func (config *Config) WriteConfig(configText []byte) {
	var out bytes.Buffer
	if len(configText) == 0 {
		configText = []byte("{}")
	}
	json.Indent(&out, configText, "", "  ")
	fileLock.Lock()
	ioutil.WriteFile(ConfigFilePath, out.Bytes(), 0755)
	fileLock.Unlock()
	config.ReadConfig()
}
