package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	_ "path/filepath"
	"time"
)

var SC ServerConfig
var Configs Config
var YamlFile []byte

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

type DataBaseConfig struct {
	DBFilePath string        `yaml:"db_file_path"`
	Name       string        `yaml:"name"`
	Type       string        `yaml:"type"`
	SlowSql    time.Duration `yaml:"slow_sql"`
}
type kafkaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Name     string `yaml:"name"`
}

type Config struct {
	DataBase []DataBaseConfig `yaml:"dataSource"`
	Redis    []RedisConfig    `yaml:"redis"`
	Kafka    kafkaConfig      `yaml:"kafka"`
	Oss      OssConfig        `yaml:"oss"`
}

type OssConfig struct {
	OssBucket       string `yaml:"ossBucket"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	OssEndPoint     string `yaml:"ossEndPoint"`
}

func InitAllConfig(fileName string) Config {
	var err error
	YamlFile, err = ioutil.ReadFile(fileName)
	if err != nil {
		panic("load config error")
	}
	dbc := Config{}
	err = yaml.Unmarshal(YamlFile, &dbc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return dbc
}

func LoadCustomizeConfig(config interface{}) error {
	err := yaml.Unmarshal(YamlFile, config)
	if err != nil {
		return err
	}
	return nil
}

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
