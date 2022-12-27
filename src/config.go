package src

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type HTTPSConfig struct {
	CertFile string
	KeyFile  string
}

type Config struct {
	Service  map[string]string `json:"service"`
	Redirect map[string]string `json:"redirect"`
	Static   map[string]string `json:"static"`
	CDN      []string          `json:"cdn"`
	HTTPS    *HTTPSConfig      `json:"HTTPS"`
}

var fileLock = &sync.RWMutex{}

var ConfigFilePath = "./config.json"

func (config *Config) Get() *Config {
	fileLock.RLock()
	defer fileLock.RUnlock()
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
		return nil
	}
	return config
}

func (config *Config) Write(data []byte) {
	var out bytes.Buffer
	if len(data) == 0 {
		data = []byte("{}")
	}
	err := json.Indent(&out, data, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fileLock.Lock()
	err = os.WriteFile(ConfigFilePath, out.Bytes(), 0755)
	if err != nil {
		log.Println(err)
		return
	}
	fileLock.Unlock()
	config.Get()
}

func (config *Config) Set(c *Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	config.Write(data)
	return err
}
