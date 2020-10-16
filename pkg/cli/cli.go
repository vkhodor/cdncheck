package cli

import (
	"flag"
	"github.com/vkhodor/cdncheck/pkg/config"
	"io/ioutil"
)

type CLIFlags struct {
	SetNormal   bool
	SetFallback bool
	GetState    bool
	Debug       bool
	ConfigFile  string
	Force       bool
	Check       bool
}

func GetArgs() CLIFlags {
	flagSetNormal := flag.Bool("set.normal", false, "set CDN in normal state")
	flagSetFallback := flag.Bool("set.fallback", false, "set CDN to fallback state without any checks")
	flagGetState := flag.Bool("get.state", false, "get CDN current state and exit")
	flagDebug := flag.Bool("debug", false, "debug mode")
	flagConfigFile := flag.String("config", "/etc/cdncheck/config.yml", "config file")
	flagForce := flag.Bool("force", false, "forcing state")
	flagCheck := flag.Bool("check", false, "check only")

	flag.Parse()

	return CLIFlags{
		SetNormal:   *flagSetNormal,
		SetFallback: *flagSetFallback,
		GetState:    *flagGetState,
		Debug:       *flagDebug,
		ConfigFile:  *flagConfigFile,
		Force:       *flagForce,
		Check:       *flagCheck,
	}
}

func GetConfig(flags *CLIFlags) (*config.YAMLConfig, error) {
	yamlFile, err := ioutil.ReadFile(flags.ConfigFile)
	if err != nil {
		return nil, err
	}
	var conf *config.YAMLConfig
	conf, err = config.NewYAMLConfig(yamlFile)
	if err != nil {
		return nil, err
	}

	if flags.Debug {
		conf.Debug = true
	}

	return conf, nil
}
