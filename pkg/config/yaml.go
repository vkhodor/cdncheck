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
	Logger   *logrus.Logger
	Debug    bool `yaml:"debug"`
	AutoBack bool `yaml:"autoback"`

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

	PolicyBasedNormal struct {
		TTL                  *int    `yaml:"ttl"`
		TrafficPolicyId      *string `yaml:"trafficPolicyId"`
		TrafficPolicyVersion *int64  `yaml:"trafficPolicyVersion"`
	}

	PolicyBasedFallback struct {
		TTL                  *int    `yaml:"ttl"`
		TrafficPolicyId      *string `yaml:"trafficPolicyId"`
		TrafficPolicyVersion *int64  `yaml:"trafficPolicyVersion"`
	}

	Normal []struct {
		Identifier           *string   `yaml:"identifier"`
		Values               *[]string `yaml:"values"`
		Type                 *string   `yaml:"type"`
		TTL                  *int      `yaml:"ttl"`
		CountryCode          *string   `yaml:"countryCode"`
		ContinentCode        *string   `yaml:"continentCode"`
		TrafficPolicyId      *string   `yaml:"trafficPolicyId"`
		TrafficPolicyVersion *int64    `yaml:"trafficPolicyVersion"`
	}

	Fallback []struct {
		Identifier           *string   `yaml:"identifier"`
		Values               *[]string `yaml:"values"`
		Type                 *string   `yaml:"type"`
		TTL                  *int      `yaml:"ttl"`
		CountryCode          *string   `yaml:"countryCode"`
		ContinentCode        *string   `yaml:"continentCode"`
		TrafficPolicyId      *string   `yaml:"trafficPolicyId"`
		TrafficPolicyVersion *int64    `yaml:"trafficPolicyVersion"`
	}

	Checks []struct {
		Name           string        `yaml:"name"`
		Domains        []string      `yaml:"domains"`
		Schema         string        `yaml:"schema"`
		Host           string        `yaml:"host"`
		Port           int           `yaml:"port"`
		Code           int           `yaml:"code"`
		Path           string        `yaml:"path"`
		TimeoutSeconds time.Duration `yaml:"timeout"`
		Retries        int           `yaml:"retries"`
		Fails          int           `yaml:"fails"`
	}
}

func (y *YAMLConfig) GetChecks() ([]checks.Check, error) {
	var chks []checks.Check
	for _, check := range y.Checks {
		switch name := check.Name; name {
		case "ssl":
			if check.Retries == 0 {
				check.Retries = 1
			}
			if check.Fails == 0 {
				check.Fails = 1
			}
			chks = append(chks,
				&checks.SSLCheck{
					Port:           check.Port,
					CertDomains:    check.Domains,
					TimeoutSeconds: check.TimeoutSeconds,
					Logger:         y.Logger,
					Retries:        check.Retries,
					Fails:          check.Fails,
				},
			)
		case "url":
			if check.Retries == 0 {
				check.Retries = 1
			}
			if check.Fails == 0 {
				check.Fails = 1
			}
			chks = append(chks,
				&checks.URLCheck{
					TimeoutSeconds: check.TimeoutSeconds,
					Port:           check.Port,
					Schema:         check.Schema,
					Path:           check.Path,
					RightCode:      check.Code,
					Logger:         y.Logger,
					Retries:        check.Retries,
					Fails:          check.Fails,
				},
			)
		default:
			return nil, errors.New("unknown check name")
		}
	}
	return chks, nil
}

func (y *YAMLConfig) GetPolicyBasedFallbackRecord() (DNSRecord, error) {
	return DNSRecord{
		TrafficPolicyId:      y.PolicyBasedFallback.TrafficPolicyId,
		TrafficPolicyVersion: y.PolicyBasedFallback.TrafficPolicyVersion,
		TTL:                  y.PolicyBasedFallback.TTL,
	}, nil
}

func (y *YAMLConfig) GetPolicyBasedNormalRecord() (DNSRecord, error) {
	return DNSRecord{
		TrafficPolicyId:      y.PolicyBasedNormal.TrafficPolicyId,
		TrafficPolicyVersion: y.PolicyBasedNormal.TrafficPolicyVersion,
		TTL:                  y.PolicyBasedNormal.TTL,
	}, nil
}

func (y *YAMLConfig) GetFallbackRecords() ([]DNSRecord, error) {
	var records []DNSRecord
	for _, r := range y.Fallback {
		if r.Identifier == nil {
			emptyString := ""
			r.Identifier = &emptyString
		}
		if r.TTL == nil {
			emptyTTL := 0
			r.TTL = &emptyTTL
		}
		if r.Values == nil {
			var emptyValues []string
			r.Values = &emptyValues
		}
		ident := fmt.Sprintf("%v:%v", *y.FallbackPrefix, *r.Identifier)
		record := DNSRecord{
			Name:                 y.Route53.RecordName,
			Values:               r.Values,
			Type:                 r.Type,
			TTL:                  r.TTL,
			CountryCode:          r.CountryCode,
			ContinentCode:        r.ContinentCode,
			Identifier:           &ident,
			TrafficPolicyId:      r.TrafficPolicyId,
			TrafficPolicyVersion: r.TrafficPolicyVersion,
		}
		records = append(records, record)
	}
	return records, nil
}

func (y *YAMLConfig) GetNormalRecords() ([]DNSRecord, error) {
	var records []DNSRecord
	for _, r := range y.Normal {
		if r.Identifier == nil {
			emptyString := ""
			r.Identifier = &emptyString
		}
		if r.TTL == nil {
			emptyTTL := 0
			r.TTL = &emptyTTL
		}
		if r.Values == nil {
			var emptyValues []string
			r.Values = &emptyValues
		}
		ident := fmt.Sprintf("%v:%v", *y.NormalPrefix, *r.Identifier)
		record := DNSRecord{
			Name:                 y.Route53.RecordName,
			Values:               r.Values,
			Type:                 r.Type,
			TTL:                  r.TTL,
			CountryCode:          r.CountryCode,
			ContinentCode:        r.ContinentCode,
			Identifier:           &ident,
			TrafficPolicyId:      r.TrafficPolicyId,
			TrafficPolicyVersion: r.TrafficPolicyVersion,
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
