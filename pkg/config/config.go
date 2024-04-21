package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Url           string `yaml:"source_url"`
	Database      string `yaml:"db_file"`
	ClientLogFile string `yaml:"client_log_file"`
	Goroutines    int    `yaml:"parallel"`
	CacheFile     string `yaml:"cache_file"`
	IndexFile     string `yaml:"index_file"`
}

func (c *Config) Init(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		fmt.Println("Чтение прошло неудачно", err)
	}
}
