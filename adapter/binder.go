package adapter

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cppforlife/go-patch/patch"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"gopkg.in/yaml.v2"

	"github.com/s-matyukevich/template-service-adapter/config"
)

type Binder struct {
	Config         *config.Config
	manifestYaml   interface{}
	deploymentYaml interface{}
}

func (b Binder) CreateBinding(bindingID string, deploymentTopology bosh.BoshVMs, manifest bosh.BoshManifest, requestParams serviceadapter.RequestParameters) (serviceadapter.Binding, error) {
	var err error
	b.manifestYaml, err = b.convert(manifest)
	if err != nil {
		return serviceadapter.Binding{}, err
	}
	b.deploymentYaml, err = b.convert(deploymentTopology)
	if err != nil {
		return serviceadapter.Binding{}, err
	}

	tmpl := template.New("binder-template")
	tmpl.Funcs(template.FuncMap{"getFromManifest": b.getTemplateFunc(b.manifestYaml)})
	tmpl.Funcs(template.FuncMap{"getFromDeployment": b.getTemplateFunc(b.deploymentYaml)})
	tmpl, err = tmpl.Parse(b.Config.BinderTemplate)
	if err != nil {
		return serviceadapter.Binding{}, err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return serviceadapter.Binding{}, err
	}

	res := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(buf.String()), &res)
	if err != nil {
		return serviceadapter.Binding{}, err
	}
	return serviceadapter.Binding{
		Credentials: res,
	}, nil
}

func (b Binder) DeleteBinding(bindingID string, deploymentTopology bosh.BoshVMs, manifest bosh.BoshManifest, requestParams serviceadapter.RequestParameters) error {
	return nil
}

func (b Binder) getTemplateFunc(doc interface{}) func(string) (string, error) {
	return func(path string) (string, error) {
		p, err := patch.NewPointerFromString(path)
		if err != nil {
			return "", err
		}
		res, err := patch.FindOp{Path: p}.Apply(doc)
		return fmt.Sprintf("%v", res), err
	}
}

func (b Binder) convert(obj interface{}) (interface{}, error) {
	doc, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var res interface{}
	err = yaml.Unmarshal(doc, &res)
	return res, err
}
