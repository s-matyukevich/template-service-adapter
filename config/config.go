package config

type Config struct {
	ManifestTemplates     map[string]string `yaml:"manifest_templates"`
	BinderTemplate        string            `yaml:"binder_template"`
	PreManifestGeneration string            `yaml:"pre_manifest_generation_script"`
	PreBinding            string            `yaml:"pre_binding_script"`
}
