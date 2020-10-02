package checks

type Check interface {
	Check(string) (state bool, err error)
}
