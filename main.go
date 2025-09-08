package main

import (
	"fmt"
	"os"

	"github.com/akornatskyy/backscribe/builders"
	"github.com/akornatskyy/backscribe/cmd/cli"
	"github.com/akornatskyy/backscribe/config"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing arguments:", err)
		os.Exit(1)
	}

	configFile, err := config.ResolveConfigFile(opts.ConfigFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error locating config file:", err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config:", err)
		os.Exit(1)
	}

	script := builders.BuildScript(cfg, configFile)
	fmt.Println(script)
}
