package utils

import (
	"gopkg.in/yaml.v2"
)

func ConvertToYamlCompatibleObject(obj interface{}) (interface{}, error) {
	doc, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var res interface{}
	err = yaml.Unmarshal(doc, &res)
	return res, err
}
