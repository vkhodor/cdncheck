package cli

import "flag"

type CLIFlags struct {
	SetNormal   bool
	SetFallback bool
	GetState    bool
	Debug       bool
}

func GetArgs() CLIFlags {
	flagSetNormal := flag.Bool("set.normal", false, "set CDN in normal state")
	flagSetFallback := flag.Bool("set.fallback", false, "set CDN to fallback state without any checks")
	flagGetState := flag.Bool("get.state", false, "get CDN current state and exit")
	flagDebug := flag.Bool("debug", false, "debug mode")

	flag.Parse()

	return CLIFlags{
		SetNormal:   *flagSetNormal,
		SetFallback: *flagSetFallback,
		GetState:    *flagGetState,
		Debug:       *flagDebug,
	}
}
