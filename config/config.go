package config

type Config struct {
	ManifestTemplates map[string]string `yaml:"manifest_templates"`
	BinderTemplate    string            `yaml:"binder_template"`
}
