package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	TcpPort  string `yaml:"tcp-port,omitempty"`
	UdpPort  string `yaml:"udp-port,omitempty"`
	HttpPort string `yaml:"http-port,omitempty"`
	Debug    bool   `yaml:"debug,omitempty"`
}

var Conf = config{}

func init() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("Need config.yaml")
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		panic(err.Error())
	}
}
