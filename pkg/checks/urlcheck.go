package checks

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type URLCheck struct {
	Path      string
	RightCode int
	Logger    *logrus.Logger
	Host      string
	Port      int
	Schema    string
}

func (h *URLCheck) Check() (bool, error) {
	if h.Schema == "" {
		h.Schema = "http"
	}
	url := h.Schema + "://" + h.Host + ":" + strconv.Itoa(h.Port) + "/" + h.Path
	h.Logger.Debug("URLCheck: url = ", url)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	h.Logger.Debug("URLCheck: ", resp.StatusCode)

	return h.checkCode(resp.StatusCode), nil
}

func (h *URLCheck) checkCode(code int) bool {
	if code != h.RightCode {
		return false
	}
	return true
}
