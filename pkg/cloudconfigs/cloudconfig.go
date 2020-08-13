package cloudconfigs

type CloudConfig interface {
	State() (string, error)
	Fallback() (bool, error)
	Normal() (bool, error)
}
