package services

import (
	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
)

type Service struct {
	IpAddress string
	Port int
	Checks *[]checks.Check
	CloudConfigs *[]cloudconfigs.CloudConfig
}
