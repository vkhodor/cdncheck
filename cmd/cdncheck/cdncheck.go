package main

import (
	"fmt"
	"github.com/vkhodor/cdncheck/pkg/cli"
	"github.com/vkhodor/cdncheck/pkg/senders"
	"github.com/vkhodor/cdncheck/pkg/servers"
	"os"
	//	"github.com/vkhodor/cdncheck/pkg/checks"
	"github.com/vkhodor/cdncheck/pkg/cloudconfigs"
)

var version = "0.1.2"

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

	logger := conf.GetLogger()

	logger.Info("debug:", conf.Debug)
	logger.Info("zoneId: " + conf.Route53.ZoneId)
	logger.Info("recordName: " + *conf.Route53.RecordName)

	var r53client cloudconfigs.CloudConfig = cloudconfigs.NewCloudRoute53(
		conf.Route53.ZoneId,
		*conf.Route53.RecordName,
		logger,
	)

	err = r53client.LoadChanges(conf)
	if err != nil {
		panic(err)
	}

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
		if currentState == *conf.FallbackPrefix && !cliFlags.Check && !cliFlags.Force {
			logger.Info(fmt.Sprintf("Current CDN state is already %v. Do nothing", *conf.FallbackPrefix))
			os.Exit(0)
		}
		_, err := r53client.Fallback()
		if err != nil {
			logger.Fatalln(fmt.Sprintf("Can't %v cloud configuration: ", *conf.FallbackPrefix), err)
		}
		logger.Info("CDN state changed to ", *conf.FallbackPrefix)
		sender.Send(fmt.Sprintf("CDN state changed to %v!", *conf.FallbackPrefix))
		os.Exit(0)
	}

	if cliFlags.SetNormal {
		if currentState == *conf.NormalPrefix && !cliFlags.Force {
			logger.Info(fmt.Sprintf("Current CDN state is already %v. Do nothing", *conf.NormalPrefix))
			os.Exit(0)
		}
		_, err = r53client.Normal()
		if err != nil {
			logger.Fatalln(fmt.Sprintf("Can't set cloud configuration to %v state: ", *conf.NormalPrefix), err)
		}
		logger.Info("CDN state changed to ", *conf.NormalPrefix)
		sender.Send(fmt.Sprintf("CDN state changed to %v.", *conf.NormalPrefix))
		os.Exit(0)
	}

	if currentState == *conf.FallbackPrefix && !cliFlags.Check && !cliFlags.Force {
		logger.Info(fmt.Sprintf("Current CDN state is already %v. Do nothing", *conf.FallbackPrefix))
		os.Exit(0)
	}

	checksList, err := conf.GetChecks()
	if err != nil {
		panic(err)
	}

	for _, host := range conf.CDNHosts {
		logger.Info("***** ", host, " *****")
		server := servers.Server{
			Logger:       logger,
			CloudConfigs: []cloudconfigs.CloudConfig{r53client},
			Checks:       checksList,
			Host:         host,
		}

		ok, err := server.Check()
		logger.Info("Check result: ", ok, err)
		if !ok {
			if cliFlags.Check {
				fmt.Println("result: ", ok)
				break
			}
			_ = sender.Send("CDN check returned error. Going to Fallback...")
			ok, err = r53client.Fallback()
			if !ok {
				logger.Fatalln(fmt.Sprintf("Can't change CDN state to %v: ", *conf.FallbackPrefix), err)
				_ = sender.Send(fmt.Sprintf("Can't change CDN state to %v: %v", *conf.FallbackPrefix, err))
			}
			logger.Info("CDN state changed to ", *conf.FallbackPrefix)
			_ = sender.Send(fmt.Sprintf("CDN state changed to %v!", *conf.FallbackPrefix))
			break
		}
	}
}
