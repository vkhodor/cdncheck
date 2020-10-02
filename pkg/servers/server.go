package servers

import (
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
)

type Server struct {
	Checks       []checks.Check
	CloudConfigs []cloudconfigs.CloudConfig
	Logger       *logrus.Logger
	Host         string
}

func (s *Server) Check() (bool, error) {
	for _, check := range s.Checks {
		ok, err := check.Check(s.Host)
		s.Logger.Debug("Check: ", ok, err)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
