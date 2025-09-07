package cli

import (
	"flag"
	"fmt"
)

type Options struct {
	ConfigFile string
	ShowHelp   bool
}

func ParseArgs() (Options, error) {
	var opts Options

	flag.StringVar(&opts.ConfigFile, "config", "", "Optional path to config file")
	flag.StringVar(&opts.ConfigFile, "c", "", "Optional path to config file (shorthand)")
	flag.BoolVar(&opts.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&opts.ShowHelp, "h", false, "Show help message (shorthand)")
	flag.Parse()

	return opts, nil
}

func PrintUsage() {
	fmt.Println(`Usage:
  backscribe [options]

Description:
  backscribe generates a sequence of shell commands based on a configuration.
  If no configuration file is specified, it searches upward from the current
  directory for a file named '[.]backscribe.(yaml|json)'.
  Also checks $HOME and $HOME/.config

Options:
  -c, --config <file>   Optional path to a config file
  -h, --help            Show this help message and exit

Examples:
  backscribe -c ./config/backscribe.yaml | sh
  backscribe | sh  # Auto-searches for config`)
}
