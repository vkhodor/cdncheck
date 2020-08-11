package cloudconfig

type CloudConfig interface {
	Status() string
	Fallback() bool
	Normal() bool
}
