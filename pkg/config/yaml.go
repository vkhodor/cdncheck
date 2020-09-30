package config

import "gopkg.in/yaml.v2"

type YAMLConfig struct {
	Debug bool `yaml:"debug"`

	Slack struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Channel  string `yaml:"channel"`
	}

	Route53 struct {
		ZoneId     string `yaml:"zoneId"`
		RecordName string `yaml:"recordName"`
	}

	CDNHosts []string `yaml:"cdnHosts"`

	Normal []struct {
		Identifier  string   `yaml:"identifier"`
		Values      []string `yaml:"values"`
		Type        string   `yaml:"type"`
		TTL         int      `yaml:"ttl"`
		CountryCode string   `yaml:"countryCode"`
	}

	Fallback []struct {
		Identifier  string   `yaml:"identifier"`
		Values      []string `yaml:"values"`
		Type        string   `yaml:"type"`
		TTL         int      `yaml:"ttl"`
		CountryCode string   `yaml:"countryCode"`
	}
}

func NewYAMLConfig(yamlData []byte) (*YAMLConfig, error) {
	var config YAMLConfig
	err := yaml.Unmarshal(yamlData, &config)
	return &config, err
}
