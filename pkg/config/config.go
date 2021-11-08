package config

import (
	"github.com/vkhodor/cdncheck/pkg/checks"
)

type Config interface {
	GetChecks() ([]checks.Check, error)
	GetFallbackRecords() ([]DNSRecord, error)
	GetNormalRecords() ([]DNSRecord, error)
	GetPolicyBasedNormalRecord() (DNSRecord, error)
	GetPolicyBasedFallbackRecord() (DNSRecord, error)
}

type DNSRecord struct {
	Name                 *string
	Identifier           *string
	Values               *[]string
	Type                 *string
	TTL                  *int
	CountryCode          *string
	ContinentCode        *string
	TrafficPolicyId      *string
	TrafficPolicyVersion *int64
}

type Check struct {
	Name    *string
	Domains *[]string
	Schema  *string
	Host    *string
	Port    *int
	Code    *int
	Path    *string
	Retries *int
	Fails   *int
}
