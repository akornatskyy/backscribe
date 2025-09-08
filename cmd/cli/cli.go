package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

type Options struct {
	ConfigFile  string
	ShowHelp    bool
	ShowVersion bool
}

func ParseArgs() (Options, error) {
	var opts Options

	flag.StringVar(&opts.ConfigFile, "config", "", "Optional path to config file")
	flag.StringVar(&opts.ConfigFile, "c", "", "Optional path to config file (shorthand)")
	flag.BoolVar(&opts.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&opts.ShowHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&opts.ShowVersion, "version", false, "Show version info and exit")
	flag.BoolVar(&opts.ShowVersion, "v", false, "Show version info and exit (shorthand)")
	flag.Parse()

	if opts.ShowHelp {
		printUsage()
		os.Exit(0)
	}

	if opts.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	return opts, nil
}

func printUsage() {
	fmt.Println(`Usage:
  backscribe [options]

Description:
  backscribe generates a sequence of shell commands based on a configuration.
  If no configuration file is specified, it searches upward from the current
  directory for a file named '[.]backscribe.(yaml|json)'.
  Also checks $HOME and $HOME/.config

Options:
  -c, --config <file>   Optional path to a config file
  -h, --help            Show this help message
	-v, --version         Show version

Examples:
  backscribe -c ./config/backscribe.yaml | sh
  backscribe | sh  # Auto-searches for config`)
}

func printVersion() {
	fmt.Printf(`Version:           %s
Go version:        %s
Git commit:        %s
Built:             %s
OS/Arch:           %s/%s
`,
		Version,
		runtime.Version(),
		GitCommit,
		BuildDate,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
