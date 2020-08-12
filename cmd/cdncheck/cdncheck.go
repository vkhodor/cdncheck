package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
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

	r53client := cloudconfigs.NewCloudRoute53("Z2WXU28CDS7KHT", "content.algorithmic.bid.", logger)

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
