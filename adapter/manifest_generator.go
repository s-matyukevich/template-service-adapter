package adapter

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	text "github.com/tonnerre/golang-text"
	"gopkg.in/yaml.v2"

	"github.com/s-matyukevich/template-service-adapter/config"
)

type ManifestGenerator struct {
	Config *config.Config
}

var GenPassword = func() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m ManifestGenerator) GenerateManifest(
	serviceDeployment serviceadapter.ServiceDeployment,
	plan serviceadapter.Plan,
	requestParams serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {
	var planName string
	if name, ok := plan.Properties["name"]; ok {
		planName = name.(string)
	} else {
		return bosh.BoshManifest{}, errors.New("Plan don't have a name property.")
	}
	tmpl := template.New("manifest-template")
	stemcellAlias := "only-stemcell"
	tmpl.Funcs(template.FuncMap{"genPassword": GenPassword})
	tmpl.Funcs(template.FuncMap{"getInstanceGroup": func(name string) (string, error) {
		for _, g := range plan.InstanceGroups {
			if g.Name == name {
				res, err := yaml.Marshal([]bosh.InstanceGroup{
					{
						Name:               g.Name,
						Instances:          g.Instances,
						VMType:             g.VMType,
						VMExtensions:       g.VMExtensions,
						PersistentDiskType: g.PersistentDiskType,
						Stemcell:           stemcellAlias,
						Networks:           m.mapNetworksToBoshNetworks(g.Networks),
						AZs:                g.AZs,
						Lifecycle:          g.Lifecycle,
					},
				})
				return string(res), err
			}
		}
		return "", fmt.Errorf("No instance group found with name %s", name)
	}})
	tmpl.Funcs(template.FuncMap{"getReleasesBlock": m.printYamlFunc("releases", m.generateReleasesBlock(serviceDeployment.Releases), "")})
	tmpl.Funcs(template.FuncMap{"getStemcellsBlock": m.printYamlFunc("stemcells", []bosh.Stemcell{{
		Alias:   stemcellAlias,
		OS:      serviceDeployment.Stemcell.OS,
		Version: serviceDeployment.Stemcell.Version,
	}}, "")})
	tmpl.Funcs(template.FuncMap{"getUpdateBlock": m.printYamlFunc("update", m.generateUpdateBlock(plan.Update), "  ")})
	_, err := tmpl.Parse(m.Config.ManifestTemplates[planName])
	if err != nil {
		return bosh.BoshManifest{}, err
	}
	b := &bytes.Buffer{}
	params := map[string]interface{}{}
	params["params"] = requestParams
	params["deployment"] = serviceDeployment
	params["plan"] = plan
	params["previousPlan"] = previousPlan
	err = tmpl.Execute(b, params)
	if err != nil {
		return bosh.BoshManifest{}, err
	}
	manifest := bosh.BoshManifest{}

	err = yaml.Unmarshal([]byte(b.String()), &manifest)
	if err != nil {
		return bosh.BoshManifest{}, err
	}
	return manifest, nil
}

func (m ManifestGenerator) printYamlFunc(blockName string, obj interface{}, indent string) func() (string, error) {
	return func() (string, error) {
		res, err := yaml.Marshal(obj)
		t := text.Indent(string(res), indent)
		return blockName + ":\n" + t, err
	}
}

func (m ManifestGenerator) generateUpdateBlock(update *serviceadapter.Update) bosh.Update {
	if update != nil {
		return bosh.Update{
			Canaries:        update.Canaries,
			MaxInFlight:     update.MaxInFlight,
			CanaryWatchTime: update.CanaryWatchTime,
			UpdateWatchTime: update.UpdateWatchTime,
			Serial:          update.Serial,
		}
	} else {
		return bosh.Update{
			Canaries:        1,
			CanaryWatchTime: "30000-240000",
			UpdateWatchTime: "30000-240000",
			MaxInFlight:     1,
		}
	}
}

func (m ManifestGenerator) generateReleasesBlock(releases serviceadapter.ServiceReleases) []bosh.Release {
	res := []bosh.Release{}
	for _, release := range releases {
		res = append(res, bosh.Release{
			Name:    release.Name,
			Version: release.Version,
		})
	}
	return res
}

func (m ManifestGenerator) mapNetworksToBoshNetworks(networks []string) []bosh.Network {
	boshNetworks := []bosh.Network{}
	for _, network := range networks {
		boshNetworks = append(boshNetworks, bosh.Network{Name: network})
	}
	return boshNetworks
}