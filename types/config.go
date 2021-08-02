package types

type ConfigFile struct {
	ClientId        string   `yaml:"clientId"`
	Tag             string   `yaml:"tag"`
	Port            int      `yaml:"port"`
	CollectInterval int      `yaml:"collectInterval"`
	SendInterval    int      `yaml:"sendInterval"`
	Servers         []Server `yaml:"servers"`
}

type Server struct {
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
}
