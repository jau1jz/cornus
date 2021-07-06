package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var SC ServerConfig
var Configs Config
var YamlFile []byte

/**
 * server config
 */
type ServerBaseConfig struct {
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"loglevel"`
	Profile  string `yaml:"profile"`
	LogPath  string `yaml:"logPath"`
	LogName  string `yaml:"logName"`
}
type ServerConfig struct {
	SConfigure ServerBaseConfig `yaml:"server"`
}

func init() {
	yamlFile, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		panic(fmt.Errorf("load application.yaml error, will exit,please fix the application"))
	}
	err = yaml.Unmarshal(yamlFile, &SC)
	if err != nil {
		panic(err)
	}
	if len(SC.SConfigure.Profile) == 0 {
		// load dev profile application-dev.yaml
		Configs = InitAllConfig("application-dev.yaml")
	} else {
		Configs = InitAllConfig(fmt.Sprintf("application-%s.yaml", SC.SConfigure.Profile))
	}
}
