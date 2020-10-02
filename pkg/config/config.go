package config

import "github.com/vkhodor/cdncheck/pkg/checks"

type Config interface {
	IsDebug() bool
	GetChecks() ([]checks.Check, error)
}
