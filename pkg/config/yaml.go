package config

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
	"gopkg.in/yaml.v2"
)

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

	Checks []struct {
		Name    string   `yaml:"name"`
		Domains []string `yaml:"domains"`
		Schema  string   `yaml:"schema"`
		Host    string   `yaml:"host"`
		Port    int      `yaml:"port"`
		Code    int      `yaml:"code"`
		Path    string   `yaml:"path"`
	}
}

func (y *YAMLConfig) GetChecks(logger *logrus.Logger) ([]checks.Check, error) {
	var chks []checks.Check
	for _, check := range y.Checks {
		switch name := check.Name; name {
		case "ssl":
			chks = append(chks,
				&checks.SSLCheck{
					Port:        check.Port,
					CertDomains: check.Domains,
					Logger:      logger,
				},
			)
		case "url":
			chks = append(chks,
				&checks.URLCheck{
					Port:      check.Port,
					Schema:    check.Schema,
					Path:      check.Path,
					RightCode: check.Code,
					Logger:    logger,
				},
			)
		default:
			return nil, errors.New("unknown check name")
		}
	}
	return chks, nil
}

func (y *YAMLConfig) GetFallbackRecords() ([]DNSRecord, error) {
	var records []DNSRecord
	for _, r := range y.Fallback {
		records = append(records, DNSRecord{
			Values:      &r.Values,
			Type:        &r.Type,
			TTL:         &r.TTL,
			CountryCode: &r.CountryCode,
			Identifier:  &r.Identifier,
		},
		)
	}
	return records, nil
}

func (y *YAMLConfig) GetNormalRecords() ([]DNSRecord, error) {
	var records []DNSRecord
	for _, r := range y.Normal {
		records = append(records, DNSRecord{
			Values:      &r.Values,
			Type:        &r.Type,
			TTL:         &r.TTL,
			CountryCode: &r.CountryCode,
			Identifier:  &r.Identifier,
		},
		)
	}
	return records, nil
}

func NewYAMLConfig(yamlData []byte) (*YAMLConfig, error) {
	var config YAMLConfig
	err := yaml.Unmarshal(yamlData, &config)
	return &config, err
}
