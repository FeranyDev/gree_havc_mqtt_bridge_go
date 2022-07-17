package config

import (
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mqtt Mqtt   `yaml:"mqtt"`
	Gree []Gree `yaml:"gree"`
}

type Mqtt struct {
	Havc  Havc  `yaml:"havc"`
	Bemfa Bemfa `yaml:"bemfa"`
}

type Havc struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ClientID string `yaml:"client_id"`
	Retain   bool   `yaml:"retain"`
	Tls      bool   `yaml:"tls"`
}
type Bemfa struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ClientID string `yaml:"client_id"`
	Retain   bool   `yaml:"retain"`
	Tls      bool   `yaml:"tls"`
}
type Gree struct {
	Host       net.IP `yaml:"host"`
	Port       int    `yaml:"port"`
	HavcTopic  string `yaml:"havc_topic"`
	BemfaTopic string `yaml:"bemfa_topic"`
}

func GetConfig(path string) *Config {
	c := &Config{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		panic(err)
	}

	return c
}
