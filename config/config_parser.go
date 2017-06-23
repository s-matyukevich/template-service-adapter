package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseConfig(configPath string) (*Config, error) {
	c := &Config{}
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, c)

	for key, val := range c.ManifestTemplates {
		t, err := ioutil.ReadFile(val)
		if err != nil {
			return nil, err
		}
		c.ManifestTemplates[key] = string(t)
	}

	t, err := ioutil.ReadFile(c.BinderTemplate)
	if err != nil {
		return nil, err
	}
	c.BinderTemplate = string(t)
	return c, err
}
