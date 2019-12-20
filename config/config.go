package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config - represent yaml configuration of loadbalancer
type Config struct {
	Listeners []int               `yaml:"listeners"`
	Targets   []ConfigTargetGroup `yaml:"targetGroups"`
}

//ConfigTargetGrup - represent a group of instances
type ConfigTargetGroup struct {
	FromPort  int      `yaml:"fromPort"`
	ToPort    int      `yaml:"toPort"`
	Name      string   `yaml:"name"`
	Path      string   `yaml:"path"`
	Sticky    bool     `yaml:"sticky"`
	Instances []string `yaml:"instances"`
}

func (ctg *ConfigTargetGroup) GetFromPort() int {
	return ctg.FromPort
}

func (ctg *ConfigTargetGroup) GetToPort() int {
	return ctg.ToPort
}

func (ctg *ConfigTargetGroup) GetPath() string {
	return ctg.Path
}

func (ctg *ConfigTargetGroup) GetInstances() []string {
	return ctg.Instances
}
func (ctg *ConfigTargetGroup) GetStickySession() bool {
	return ctg.Sticky
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
