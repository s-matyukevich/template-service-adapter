package config

type Config struct {
	ManifestTemplates map[string]string
	BinderTemplate    string
	ServiceName       string
}
