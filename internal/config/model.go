package config

type Config struct {
	Hosts []Host `yaml:"hosts,omitempty"`
}

type Host struct {
	Host  string `yaml:"host"`
	OAuth *OAuth `yaml:"oauth,omitempty"`
}

type OAuth struct {
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	TokenUrl     string `yaml:"tokenUrl"`
}
