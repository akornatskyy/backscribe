package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Archive struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"`
	Src      []string `json:"src" yaml:"src"`
	Dst      string   `json:"dst,omitempty" yaml:"dst,omitempty"`
	Method   *Method  `json:"method,omitempty" yaml:"method,omitempty"`
	Files    []string `json:"files" yaml:"files"`
	Exclude  []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
	Rexclude []string `json:"rexclude,omitempty" yaml:"rexclude,omitempty"`
	Copy     []string `json:"copy,omitempty" yaml:"copy,omitempty"`
}

type Method struct {
	Level *int `json:"level,omitempty" yaml:"level,omitempty"`
}

type Group struct {
	Name     string    `json:"name" yaml:"name"`
	Archives []Archive `json:"archives" yaml:"archives"`
	Skip     bool      `json:"skip,omitempty" yaml:"skip,omitempty"`
}

type Config struct {
	Groups []Group `json:"groups" yaml:"groups"`
}

func quote(name string) string {
	if (strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`)) ||
		(strings.HasPrefix(name, `'`) && strings.HasSuffix(name, `'`)) {
		return name
	}
	if strings.Contains(name, `\ `) {
		return name
	}
	if !strings.Contains(name, " ") {
		return name
	}
	if strings.HasPrefix(name, "~") {
		return strings.ReplaceAll(name, ` `, `\ `)
	}
	return `"` + strings.ReplaceAll(name, `"`, `\"`) + `"`
}

func buildArchive(archive Archive, group string) string {
	var src []string
	t := archive.Type
	if t == "cp" {
		target := ""
		if archive.Dst != "" {
			target = "/" + archive.Dst
			src = append(src, fmt.Sprintf("\n  mkdir -p \"${DEST_DIR}%s\"", target))
		}

		files := make([]string, len(archive.Src))
		for i, f := range archive.Src {
			files[i] = quote(f)
		}

		src = append(src, fmt.Sprintf(`
  log '\e[32m▶\e[0m %s => %s'
  cp -a -n %s "${DEST_DIR}%s"
`, group, archive.Name, strings.Join(archive.Src, " "), target))
		return strings.Join(src, "")
	}

	options := []string{"-t" + t, "-bso0"}
	if t == "7z" && archive.Method != nil && archive.Method.Level != nil {
		options = append(options, fmt.Sprintf("-mx%d", *archive.Method.Level))
	}

	src = append(src, fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.%s" ]; then
    log '\e[32m▶\e[0m %s => %s'
    7z a %s "${DEST_DIR}/%s.%s" \
      `,
		archive.Name, t, group, archive.Name,
		strings.Join(options, " "), archive.Name, t,
	))

	filesAndExcludes := []string{}
	for _, f := range archive.Files {
		filesAndExcludes = append(filesAndExcludes, quote(f))
	}
	if archive.Exclude != nil {
		for _, x := range archive.Exclude {
			filesAndExcludes = append(filesAndExcludes, "-x!"+quote(x))
		}
	}
	if archive.Rexclude != nil {
		for _, x := range archive.Rexclude {
			filesAndExcludes = append(filesAndExcludes, "-xr!"+quote(x))
		}
	}
	if archive.Copy != nil {
		for _, x := range archive.Copy {
			filesAndExcludes = append(filesAndExcludes, "-xr!"+quote(x))
		}
	}

	src = append(src, strings.Join(filesAndExcludes, " \\\n      "))

	if t == "7z" && archive.Copy != nil {
		options = filterOut(options, func(o string) bool {
			return !strings.HasPrefix(o, "-mx")
		})
		options = append(options, "-mx0")

		src = append(src, fmt.Sprintf(`

    7z u %s "${DEST_DIR}/%s.%s" \
      `, strings.Join(options, " "), archive.Name, t))

		filesAndExcludes = []string{}
		for _, f := range archive.Files {
			filesAndExcludes = append(filesAndExcludes, quote(f))
		}
		if archive.Exclude != nil {
			for _, x := range archive.Exclude {
				filesAndExcludes = append(filesAndExcludes, "-x!"+quote(x))
			}
		}
		if archive.Rexclude != nil {
			for _, x := range archive.Rexclude {
				filesAndExcludes = append(filesAndExcludes, "-xr!"+quote(x))
			}
		}
		if archive.Copy != nil {
			for _, x := range archive.Copy {
				filesAndExcludes = append(filesAndExcludes, "-ir!"+quote(x))
			}
		}

		src = append(src, strings.Join(filesAndExcludes, " \\\n      "))
	}

	src = append(src, `
  else
    log '\e[33m↷\e[0m `+group+` => `+archive.Name+`'
  fi
`)

	return strings.Join(src, "")
}

func filterOut(slice []string, keep func(string) bool) []string {
	var result []string
	for _, s := range slice {
		if keep(s) {
			result = append(result, s)
		}
	}
	return result
}

func buildGroup(group Group) string {
	src := []string{fmt.Sprintf("\nbackup_%s() {", group.Name)}
	for _, archive := range group.Archives {
		src = append(src, buildArchive(archive, group.Name))
	}
	src = append(src, "}")
	return strings.Join(src, "")
}

func buildScript(config *Config, configFile string) string {
	src := []string{`#!/bin/sh
set -o errexit
`}
	src = append(src, fmt.Sprintf(`
CONFIG_FILE=%s
DEST_DIR=~/backups/$(date '+%%Y-%%m-%%d')
START=$(date +%%s)
`, strconv.Quote(configFile)))
	src = append(src, `
log() {
  NOW=$(date +%s)
  ELAPSED=$((NOW - START))
  MIN=$((ELAPSED / 60))
  SEC=$((ELAPSED % 60))
  if [ "${MIN}" -lt 10 ]; then MIN="0${MIN}"; fi
  if [ "${SEC}" -lt 10 ]; then SEC="0${SEC}"; fi
  echo -e "\033[90m${MIN}:${SEC}\033[0m" $1
}
`)

	for _, group := range config.Groups {
		src = append(src, buildGroup(group))
		src = append(src, "\n")
	}

	src = append(src, `
echo "CONFIG_FILE=${CONFIG_FILE}"
echo "DEST_DIR=${DEST_DIR}"
mkdir -p "${DEST_DIR}"
`)

	for _, group := range config.Groups {
		prefix := ""
		if group.Skip {
			prefix = "# "
		}
		src = append(src, fmt.Sprintf("\n%sbackup_%s", prefix, group.Name))
	}

	src = append(src, "\n\nlog '\\e[32m✓\\e[0m all done'")
	return strings.Join(src, "")
}

func findFirstMatchingFile(filenames []string) (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := startDir

	for {
		for _, filename := range filenames {
			path := filepath.Join(dir, filename)
			if fileExists(path) {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range filenames {
			candidate := filepath.Join(home, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}

		configDir := filepath.Join(home, ".config")
		for _, name := range filenames {
			candidate := filepath.Join(configDir, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("no matching config file found in any parent directory")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, cfg)
	case ".json":
		err = json.Unmarshal(data, cfg)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
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
  -h, --help            Show this help message and exit

Examples:
  backscribe -c ./config/backscribe.yaml | sh
  backscribe | sh  # Auto-searches for config`)
}

func main() {
	var configFile string
	var showHelp bool
	flag.StringVar(&configFile, "config", "", "Optional path to JSON config file")
	flag.StringVar(&configFile, "c", "", "Optional path to JSON config file (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if configFile == "" {
		var err error
		patterns := []string{
			"backscribe.yaml", "backscribe.yml", "backscribe.json",
			".backscribe.yaml", ".backscribe.yml", ".backscribe.json",
		}
		configFile, err = findFirstMatchingFile(patterns)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config:", err)
		os.Exit(1)
	}

	script := buildScript(config, configFile)
	fmt.Println(script)
}
