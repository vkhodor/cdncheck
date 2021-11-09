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
	logger.Info("autoback: ", conf.AutoBack)
	logger.Info("zoneId: " + conf.Route53.ZoneId)
	logger.Info("recordName: " + *conf.Route53.RecordName)
	var r53client cloudconfigs.CloudConfig = cloudconfigs.NewCloudRoute53(
		conf.Route53.ZoneId,
		*conf.Route53.RecordName,
		logger,
	)

	if true {
		r53client = cloudconfigs.NewCloudRoute53PolicyBased(
			conf.Route53.ZoneId,
			*conf.Route53.RecordName,
			logger,
			conf.FallbackPrefix,
			conf.NormalPrefix,
		)
	}

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
		_ = sender.Send(fmt.Sprintf("[%v] CDN state changed to %v!", *conf.Route53.RecordName, *conf.FallbackPrefix))
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
		_ = sender.Send(fmt.Sprintf("[%v] CDN state changed to %v.", *conf.Route53.RecordName, *conf.NormalPrefix))
		os.Exit(0)
	}

	if currentState == *conf.FallbackPrefix && !cliFlags.Check && !cliFlags.Force {
		logger.Info(fmt.Sprintf("Current CDN state is already %v. Do nothing", *conf.FallbackPrefix))
		logger.Debug(conf.Slack.AlwaysFallbackSend)
		if conf.Slack.AlwaysFallbackSend {
			_ = sender.Send(fmt.Sprintf("[%v] Current CDN state is %v.", *conf.Route53.RecordName, *conf.FallbackPrefix))
		}
		if !conf.AutoBack {
			os.Exit(0)
		}
		logger.Info(fmt.Sprintf("AutoBack is true. Going to check and set normal state if all checks will be okay!"))
	}

	checksList, err := conf.GetChecks()
	if err != nil {
		panic(err)
	}

	checksResult := true
	for _, host := range conf.CDNHosts {
		logger.Info("***** ", host, " *****")
		server := servers.Server{
			Logger:       logger,
			CloudConfigs: []cloudconfigs.CloudConfig{r53client},
			Checks:       checksList,
			Host:         host,
		}

		ok, err := server.Check()
		checksResult = ok
		logger.Info("Check result: ", ok, err)
		if !ok {
			if cliFlags.Check {
				fmt.Println("result: ", ok)
				break
			}
			if currentState != *conf.FallbackPrefix {
				_ = sender.Send(fmt.Sprintf("[%v] CDN check returned error. Going to Fallback...", *conf.Route53.RecordName))
				ok, err = r53client.Fallback()
				if !ok {
					logger.Fatalln(fmt.Sprintf("Can't change CDN state to %v: ", *conf.FallbackPrefix), err)
					_ = sender.Send(fmt.Sprintf("[%v] Can't change CDN state to %v: %v", *conf.Route53.RecordName, *conf.FallbackPrefix, err))
				}
				logger.Info("CDN state changed to ", *conf.FallbackPrefix)
				_ = sender.Send(fmt.Sprintf("[%v] CDN state changed to %v!", *conf.Route53.RecordName, *conf.FallbackPrefix))
			}
			break
		}
	}

	if checksResult && currentState == *conf.FallbackPrefix && !cliFlags.Check && conf.AutoBack {
		_, err = r53client.Normal()
		if err != nil {
			logger.Fatalln(fmt.Sprintf("Can't set cloud configuration to %v state: ", *conf.NormalPrefix), err)
		}
		logger.Info("CDN state changed to ", *conf.NormalPrefix)
		_ = sender.Send(fmt.Sprintf("[%v] CDN state changed to %v.", *conf.Route53.RecordName, *conf.NormalPrefix))
		os.Exit(0)
	}
}
