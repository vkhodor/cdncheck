package main

import (
	logrus "github.com/sirupsen/logrus"
	//	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
	//	"github.com/vkhodor/cdncheck/pkg/servers"
	"os"
)

var version = "0.0.1"

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

func main() {
	logger := NewLogger()

	zoneId := "Z2WXU28CDS7KHT"
	recordName := "content.algorithmic.bid."

	//zoneId := "Z075237433HUE55HSW0Z7"
	//recordName := "content.cdn.personaly.bid."

	logger.Info("zoneId: " + zoneId)
	logger.Info("recordName: " + recordName)

	var r53client cloudconfigs.CloudConfig = cloudconfigs.NewCloudRoute53(zoneId, recordName, logger)
	/*
		sslCheck := &checks.SSLCheck{
			CertDomains: []string{
				"content.cdn.personaly.bid",
				"us-01.cdn.personaly.bid",
				"*.cdn.personaly.bid",
			},
			Logger: logger,
			Host:   "us-01.cdn.personaly.bid",
			Port:   443,
		}

		httpCheck := &checks.URLCheck{
			Path:      "checks/status.txt",
			RightCode: 200,
			Logger:    logger,
			Host:      "us-01.cdn.personaly.bid",
			Port:      80,
			Schema:    "http",
		}

		httpsCheck := &checks.URLCheck{
			Path:      "checks/status.txt",
			RightCode: 200,
			Logger:    logger,
			Host:      "199.115.113.118",
			Port:      443,
			Schema:    "https",
		}
	*/
	/*
		usServer01 := servers.Server{
			Logger:       logger,
			CloudConfigs: []cloudconfigs.CloudConfig{r53client},
			Checks: []checks.Check{
				sslCheck,
				httpCheck,
				httpsCheck,
			},
		}

		usServer01.TryToFallback()
	*/

	if state, err := r53client.State(); state == "normal" {
		if err != nil {
			logger.Fatalln("can't get state of cloud configuration: ", err)
		}
		logger.Info("current cloud configuration state is normal")

		_, err = r53client.Fallback()
		if err != nil {
			logger.Fatalln("can't fallback cloud configuration: ", err)
		}
		logger.Info("cloud configuration state changed to fallback")
		os.Exit(0)
	}
	logger.Info("current cloud configuration state is fallback")
	_, err := r53client.Normal()
	if err != nil {
		logger.Fatalln("can't back cloud configuration to normal state: ", err)
	}
	logger.Info("cloud configuration state changed to normal state")
	os.Exit(0)

}
