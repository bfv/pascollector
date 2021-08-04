package types

type ConfigFile struct {
	ClientId        string        `yaml:"clientId"`
	Server          string        `yaml:"server"`
	Tag             string        `yaml:"tag"`
	Port            int           `yaml:"port"`
	CollectInterval int           `yaml:"collectInterval"`
	SendInterval    int           `yaml:"sendInterval"`
	PasInstances    []PasInstance `mapstructure:"instances"`
}

type PasInstance struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Creds string `yaml:"creds"`
}
