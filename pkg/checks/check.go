package checks

type Check interface {
	Check() (state bool, err error)
}

type HTTPSCheck struct {
}

func (h *HTTPSCheck) Check(string, int) (bool, error) {
	return true, nil
}

type CachedContentCheck struct {
}

func (h *CachedContentCheck) Check(string, int) (bool, error) {
	return true, nil
}

type NonCachedContentCheck struct {
}

func (h *NonCachedContentCheck) Check(string, int) (bool, error) {
	return true, nil
}
