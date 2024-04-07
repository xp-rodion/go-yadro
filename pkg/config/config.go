package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Url           string `yaml:"source_url"`
	Database      string `yaml:"db_file"`
	ClientLogFile string `yaml:"client_log_file"`
}

func (c *Config) Init(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(data, &c)
}
