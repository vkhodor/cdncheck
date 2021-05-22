package config

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
	"gopkg.in/yaml.v2"
	"time"
)

type YAMLConfig struct {
	Logger *logrus.Logger
	Debug  bool `yaml:"debug"`

	Slack struct {
		URL                string `yaml:"url"`
		Username           string `yaml:"username"`
		Channel            string `yaml:"channel"`
		AlwaysFallbackSend bool   `yaml:"alwaysFallbackSend"`
	}

	Route53 struct {
		ZoneId     string  `yaml:"zoneId"`
		RecordName *string `yaml:"recordName"`
	}

	CDNHosts       []string `yaml:"cdnHosts"`
	NormalPrefix   *string  `yaml:"normalPrefix"`
	FallbackPrefix *string  `yaml:"fallbackPrefix"`

	Normal []struct {
		Identifier    *string   `yaml:"identifier"`
		Values        *[]string `yaml:"values"`
		Type          *string   `yaml:"type"`
		TTL           *int      `yaml:"ttl"`
		CountryCode   *string   `yaml:"countryCode"`
		ContinentCode *string   `yaml:"continentCode"`
	}

	Fallback []struct {
		Identifier    *string   `yaml:"identifier"`
		Values        *[]string `yaml:"values"`
		Type          *string   `yaml:"type"`
		TTL           *int      `yaml:"ttl"`
		CountryCode   *string   `yaml:"countryCode"`
		ContinentCode *string   `yaml:"continentCode"`
	}

	Checks []struct {
		Name    string   `yaml:"name"`
		Domains []string `yaml:"domains"`
		Schema  string   `yaml:"schema"`
		Host    string   `yaml:"host"`
		Port    int      `yaml:"port"`
		Code    int      `yaml:"code"`
		Path    string   `yaml:"path"`
		TimeoutSeconds time.Duration `yaml:"timeout"`
	}
}

func (y *YAMLConfig) GetChecks() ([]checks.Check, error) {
	var chks []checks.Check
	for _, check := range y.Checks {
		switch name := check.Name; name {
		case "ssl":
			chks = append(chks,
				&checks.SSLCheck{
					Port:        check.Port,
					CertDomains: check.Domains,
					Logger:      y.Logger,
				},
			)
		case "url":
			chks = append(chks,
				&checks.URLCheck{
					TimeoutSeconds: check.TimeoutSeconds,
					Port:      check.Port,
					Schema:    check.Schema,
					Path:      check.Path,
					RightCode: check.Code,
					Logger:    y.Logger,
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
		ident := fmt.Sprintf("%v:%v", *y.FallbackPrefix, *r.Identifier)
		record := DNSRecord{
			Name:          y.Route53.RecordName,
			Values:        r.Values,
			Type:          r.Type,
			TTL:           r.TTL,
			CountryCode:   r.CountryCode,
			ContinentCode: r.ContinentCode,
			Identifier:    &ident,
		}
		records = append(records, record)
	}
	return records, nil
}

func (y *YAMLConfig) GetNormalRecords() ([]DNSRecord, error) {
	var records []DNSRecord
	for _, r := range y.Normal {
		ident := fmt.Sprintf("%v:%v", *y.NormalPrefix, *r.Identifier)
		record := DNSRecord{
			Name:          y.Route53.RecordName,
			Values:        r.Values,
			Type:          r.Type,
			TTL:           r.TTL,
			CountryCode:   r.CountryCode,
			ContinentCode: r.ContinentCode,
			Identifier:    &ident,
		}
		records = append(records, record)
	}
	return records, nil
}

func (y *YAMLConfig) GetLogger() *logrus.Logger {
	return y.Logger
}

func NewYAMLConfig(yamlData []byte) (*YAMLConfig, error) {
	var config YAMLConfig
	err := yaml.Unmarshal(yamlData, &config)

	level := logrus.InfoLevel
	if config.Debug {
		level = logrus.DebugLevel
	}
	config.Logger = NewLogger(level)

	return &config, err
}

func NewLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(level)
	return logger
}
