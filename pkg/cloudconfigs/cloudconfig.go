package cloudconfigs

import "github.com/vkhodor/cdncheck/pkg/config"

type CloudConfig interface {
	State() (string, error)
	Fallback() (bool, error)
	Normal() (bool, error)
	LoadChanges(config config.Config) error
}
