package lib

import (
	"bytes"
	"encoding/json"
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

type StaticConfig struct {
	Path string
	Dir  string
}

type RedirectConfig struct {
	Path   string
	Target string
}

type WebsocketConfig struct {
	Path   string
	Target string
}

type Config struct {
	Password        string
	Service         []ServiceConfig
	Redirect        []RedirectConfig
	Static          []StaticConfig
	DynamicService  *DynamicServiceConfig
	CDNList         []string
	WebsocketConfig []WebsocketConfig
	HTTPS           *HTTPSConfig
}

var fileLock = &sync.RWMutex{}

var ConfigFilePath = "./config.json"

func (config *Config) Get() *Config {
	fileLock.RLock()
	defer fileLock.RUnlock()
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&config)
	return config
}

func (config *Config) Write(data []byte) {
	var out bytes.Buffer
	if len(data) == 0 {
		data = []byte("{}")
	}
	json.Indent(&out, data, "", "  ")
	fileLock.Lock()
	os.WriteFile(ConfigFilePath, out.Bytes(), 0755)
	fileLock.Unlock()
	config.Get()
}

func (config *Config) Set(c *Config) error {
	c.Password = config.Password
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	config.Write(data)
	return err
}
