package checks

type Check interface {
	Check(ip string, port int)(state bool, err error)
}



type HTTPCheck struct {
}
func (h *HTTPCheck) Check(string, int)(bool, error) {
	return true, nil
}

type HTTPSCheck struct {
}
func (h *HTTPSCheck) Check(string, int)(bool, error) {
	return true, nil
}


type SSLCheck struct {
}
func (h *SSLCheck) Check(string, int)(bool, error) {
	return true, nil
}


type CachedContentCheck struct {
}
func (h *CachedContentCheck) Check(string, int)(bool, error) {
	return true, nil
}

type NonCachedContentCheck struct {
}
func (h *NonCachedContentCheck) Check(string, int)(bool, error) {
	return true, nil
}