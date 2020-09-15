package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cli"
	"github.com/vkhodor/cdncheck/pkg/servers"
	"os"

	//	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
)

var version = "0.0.1"

func NewLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(level)
	return logger
}

func main() {
	hosts := []string{
		"us-01.cdn.personaly.bid",
		"us-02.cdn.personaly.bid",
		"eu-01.cdn.personaly.bid",
		"jp-01.cdn.personaly.bid",
	}

	zoneId := "Z2WXU28CDS7KHT"
	recordName := "content.algorithmic.bid."

	//zoneId := "Z075237433HUE55HSW0Z7"
	//recordName := "content.cdn.personaly.bid."

	cliFlags := cli.GetArgs()
	level := logrus.InfoLevel
	if cliFlags.Debug {
		level = logrus.DebugLevel
	}
	logger := NewLogger(level)

	logger.Info("zoneId: " + zoneId)
	logger.Info("recordName: " + recordName)

	var r53client cloudconfigs.CloudConfig = cloudconfigs.NewCloudRoute53(zoneId, recordName, logger)
	currentState, err := r53client.State()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Info("Current CDN state: ", currentState)
	if cliFlags.GetState {
		os.Exit(0)
	}

	if cliFlags.SetFallback {
		if currentState == "fallback" {
			logger.Info("Current CDN state is already fallback. Do nothing")
			os.Exit(0)
		}
		_, err := r53client.Fallback()
		if err != nil {
			logger.Fatalln("Can't fallback cloud configuration: ", err)
		}
		logger.Info("CDN state changed to fallback")
		os.Exit(0)
	}

	if cliFlags.SetNormal {
		if currentState == "normal" {
			logger.Info("Current CDN state is already normal. Do nothing")
			os.Exit(0)
		}
		_, err = r53client.Normal()
		if err != nil {
			logger.Fatalln("Can't back cloud configuration to normal state: ", err)
		}
		logger.Info("CDN state changed to normal")
		os.Exit(0)
	}

	if currentState == "fallback" {
		logger.Info("Current CDN state is already fallback. Do nothing")
		os.Exit(0)
	}
	for _, host := range hosts {
		logger.Info("***** ", host, " *****")
		sslCheck := &checks.SSLCheck{
			CertDomains: []string{
				"content.cdn.personaly.bid",
				"*.cdn.personaly.bid",
				host,
			},
			Logger: logger,
			Host:   host,
			Port:   443,
		}

		httpCheck := &checks.URLCheck{
			Path:      "checks/status.txt",
			RightCode: 200,
			Logger:    logger,
			Host:      host,
			Port:      80,
			Schema:    "http",
		}

		httpsCheck := &checks.URLCheck{
			Path:      "checks/status.txt",
			RightCode: 200,
			Logger:    logger,
			Host:      host,
			Port:      443,
			Schema:    "https",
		}

		server := servers.Server{
			Logger:       logger,
			CloudConfigs: []cloudconfigs.CloudConfig{r53client},
			Checks: []checks.Check{
				sslCheck,
				httpCheck,
				httpsCheck,
			},
		}

		ok, err := server.Check()
		logger.Info("Check result: ", ok, err)
		if !ok {
			ok, err = r53client.Fallback()
			if !ok {
				logger.Fatalln("Can't change CDN state to fallback: ", err)
			}
			logger.Info("CDN state changed to fallback")
		}
	}
}
