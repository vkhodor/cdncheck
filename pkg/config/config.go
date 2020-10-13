package config

import (
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
)

type Config interface {
	GetChecks(*logrus.Logger) ([]checks.Check, error)
	GetFallbackRecords() ([]DNSRecord, error)
	GetNormalRecords() ([]DNSRecord, error)
}

type DNSRecord struct {
	Name        *string
	Identifier  *string
	Values      *[]string
	Type        *string
	TTL         *int
	CountryCode *string
}

type Check struct {
	Name    *string
	Domains *[]string
	Schema  *string
	Host    *string
	Port    *int
	Code    *int
	Path    *string
}
