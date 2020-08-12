package cloudconfigs

type CloudConfig interface {
	State() string
	Fallback() bool
	Normal() bool
}
