package config

import (
	"testing"
)

func TestNewYAMLConfig(t *testing.T) {
	yaml := `
---
debug: on

slack:
  url: 'https://hooks.slack.com/services/111'
  username: 'cdn'
  channel: 'some-channel'

route53:
  zoneId: '123'
  recordName: 'qwerty.com.'

cdnHosts:
  - 'u-01.cdn.qwery.com'
  - 'u-02.cdn.qwerty.com'
  - 'e-01.cdn.qwerty.com'
  - 'j-01.cdn.qwerty.com'

normal:
  - identifier: 'default-content'
    values:
      - '1.2.3.4'
      - '1.1.1.1'
    type: 'A'
    ttl: 60
    countryCode: '*'

  - identifier: 'u-content'
    values:
      - '127.0.0.1'
      - '127.0.0.2'
    type: 'A'
    ttl: 60
    countryCode: 'NA'

  - identifier: 'j-content'
    values:
      - '8.8.8.8'
    type: 'A'
    ttl: 60
    countryCode: 'JP'

  - identifier: 'a-content'
    values:
      - '4.4.4.4'
    type: 'A'
    ttl: 60
    countryCode: 'AS'

  - identifier: 'e-content'
    values:
      - '5.5.5.5'
    type: 'A'
    ttl: 60
    countryCode: 'EU'

fallback:
  - values:
      - 'xxxx.cloudfront.net'
    type: 'CNAME'
    ttl: 60
`

	cfg, _ := NewYAMLConfig([]byte(yaml))

	if cfg.Debug != true {
		t.Error()
	}
	if cfg.Slack.URL != "https://hooks.slack.com/services/111" {
		t.Error()
	}
	if cfg.Slack.Username != "cdn" {
		t.Error()
	}
	if cfg.Slack.Channel != "some-channel" {
		t.Error()
	}
	if cfg.Route53.ZoneId != "123" {
		t.Error()
	}
	if cfg.Route53.RecordName != "qwerty.com." {
		t.Error()
	}
	if len(cfg.CDNHosts) != 4 {
		t.Error()
	}

	flag := false
	for _,host := range cfg.CDNHosts {
			if host == "j-01.cdn.qwerty.com" {
				flag = true
			}
	}
	if !flag {
		t.Error()
	}

	if len(cfg.Normal) != 5 {
		t.Error()
	}

	for _, record := range cfg.Normal {
		if record.Identifier == "u-content" {
			if len(record.Values) != 2 {
				t.Error()
			}
			if record.TTL != 60 {
				t.Error()
			}
			if record.Type != "A" {
				t.Error()
			}
			if record.CountryCode != "NA" {
				t.Error()
			}
		}
	}

	if len(cfg.Fallback) != 1 {
		t.Error()
	}
	if cfg.Fallback[0].Values[0] != "xxxx.cloudfront.net" {
		t.Error()
	}
}
