package config

import (
	"github.com/docopt/docopt-go"
)

var usage string = `ims.

Usage:
  ims [--storage=<src>] [--cache=<cch>] [--address=<address>]
      [--log-level=<level>] [--no-optimization]
  ims -h | --help
  ims --version

Options:
  -h --help                      Show this screen.
  --version                      Show version.
  --storage=<src>                Storage backend                    [default: ./].
  --address=<address>            Address that ims should bind to    [default: :5995].
  --log-level=<level>            Log level (debug/info/warn/error)  [default: info].
  --no-optimization              Disables image optimization.`

type Config struct {
	Storage        string
	Address        string
	LogLevel       string
	NoOptimization bool
}

func ParseCLIConfig() *Config {
	// parse command line args, exiting if required
	arguments, _ := docopt.Parse(usage, nil, true, "ims 0.1", false)
	return &Config{
		Storage:        arguments["--storage"].(string),
		Address:        arguments["--address"].(string),
		LogLevel:       arguments["--log-level"].(string),
		NoOptimization: arguments["--no-optimization"].(bool),
	}
}
