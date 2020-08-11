package cloudconfig

type CloudConfig interface {
	Connect() interface{}
	Status() string
	Fallback() bool
	Normal() bool
}
