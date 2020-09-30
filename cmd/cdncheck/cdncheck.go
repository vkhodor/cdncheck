package main

import (
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cli"
	"github.com/vkhodor/cdncheck/pkg/senders"
	"github.com/vkhodor/cdncheck/pkg/servers"
	"os"
	//	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
)

var version = "0.0.2"

func NewLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(level)
	return logger
}

func main() {
	cliFlags := cli.GetArgs()
	fmt.Println(cliFlags.ConfigFile)
	conf, err := cli.GetConfig(&cliFlags)
	if err != nil {
		panic(err)
	}

	sender := senders.NewSlack(
		conf.Slack.URL,
		conf.Slack.Username,
		conf.Slack.Channel,
	)

	level := logrus.InfoLevel
	if conf.Debug {
		level = logrus.DebugLevel
	}
	logger := NewLogger(level)

	logger.Info("zoneId: " + conf.Route53.ZoneId)
	logger.Info("recordName: " + conf.Route53.RecordName)

	var r53client cloudconfigs.CloudConfig = cloudconfigs.NewCloudRoute53(
		conf.Route53.ZoneId,
		conf.Route53.RecordName,
		logger,
	)
	var currentState string
	currentState, err = r53client.State()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Info("Current CDN state: ", currentState)
	if cliFlags.GetState {
		fmt.Println(currentState)
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

	for _, host := range conf.CDNHosts {
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
			_ = sender.Send("CDN check returned error. Going to Fallback...")
			ok, err = r53client.Fallback()
			if !ok {
				logger.Fatalln("Can't change CDN state to fallback: ", err)
				_ = sender.Send(fmt.Sprintf("Can't change CDN state to fallback: %v", err))
			}
			logger.Info("CDN state changed to fallback")
			_ = sender.Send("CDN state changed to fallback!")
			break
		}
	}
}
