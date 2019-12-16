package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Listeners []int `yaml:"listeners"`
	Targets   []struct {
		FromPort  int      `yaml:"fromPort"`
		ToPort    int      `yaml:"toPort"`
		Name      string   `yaml:"name"`
		Instances []string `yaml:"instances"`
	} `yaml:"targetGroups"`
}

func InitConfig(filename string) (c *Config, err error) {
	yml, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err = parseYaml(yml)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func parseYaml(yml []byte) (c *Config, err error) {
	err = yaml.Unmarshal(yml, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
