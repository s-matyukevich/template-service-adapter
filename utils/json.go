package utils

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

func MakeJsonCompatible(manifest bosh.BoshManifest) bosh.BoshManifest {
	manifest.Properties = makeJsonCompatibleMap(manifest.Properties)
	for _, group := range manifest.InstanceGroups {
		group.Properties = makeJsonCompatibleMap(group.Properties)
		for _, job := range group.Jobs {
			job.Properties = makeJsonCompatibleMap(job.Properties)
		}
	}
	return manifest
}

func makeJsonCompatibleMap(obj map[string]interface{}) map[string]interface{} {
	for k, v := range obj {
		obj[k] = convert(v)
	}
	return obj
}
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
