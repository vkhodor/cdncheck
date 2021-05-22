package checks

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type URLCheck struct {
	Path      string
	RightCode int
	Logger    *logrus.Logger
	Port      int
	Schema    string
	TimeoutSeconds time.Duration
}

func (h *URLCheck) Check(host string) (bool, error) {
	if h.Schema == "" {
		h.Schema = "http"
	}
	url := h.Schema + "://" + host + ":" + strconv.Itoa(h.Port) + "/" + h.Path
	h.Logger.Debug("URLCheck: url = ", url)

	client := http.Client{
		Timeout: h.TimeoutSeconds * time.Second,
	}

	resp, err := client.Get(url)
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
