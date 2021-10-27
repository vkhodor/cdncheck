package checks

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type URLCheck struct {
	Path           string
	RightCode      int
	Logger         *logrus.Logger
	Port           int
	Schema         string
	TimeoutSeconds time.Duration
	Retries        int
	Fails          int
}

func (h *URLCheck) Check(host string) (bool, error) {

	h.Logger.Debug("Retries: ", h.Retries)
	h.Logger.Debug("Fails: ", h.Fails)
	fails := 0
	for i := 0; i < h.Retries; i++ {
		if fails >= h.Fails {
			return false, nil
		}
		h.Logger.Debug("Retry: ", i+1)

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
			fails += 1
		}

		h.Logger.Debug("URLCheck: ", resp.StatusCode)

		if !h.checkCode(resp.StatusCode) {
			fails += 1
		}
	}
	if fails >= h.Fails {
		return false, nil
	}
	return true, nil
}

func (h *URLCheck) checkCode(code int) bool {
	if code != h.RightCode {
		return false
	}
	return true
}
