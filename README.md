# backscribe

[![tests](https://github.com/akornatskyy/backscribe/actions/workflows/tests.yml/badge.svg)](https://github.com/akornatskyy/backscribe/actions/workflows/tests.yml)

*Backscribe* is a flexible command generator for archiving, copying, and backing up files using a structured configuration definitions. It outputs a series of shell commands that you can review or pipe directly into your shell (e.g., via `sh` or `bash`).

---

## Usage

```text
Usage:
  backscribe [options]

Description:
  backscribe generates a sequence of shell commands based on a configuration.
  If no configuration file is specified, it searches upward from the current directory
  for a file named '[.]backscribe.(yaml|json)'.
  Also checks $HOME and $HOME/.config

Options:
  -c, --config <file>   Optional path to a config file
  -h, --help            Show this help message and exit

Examples:
  backscribe -c ./config/backscribe.yaml | sh
  backscribe | sh  # Auto-searches for config
```

## Prerequisites

The 7z ([7-Zip](https://sourceforge.net/projects/sevenzip/files/7-Zip/)) command-line tool must be installed and discoverable in your system’s PATH. You should be able to run 7z from any terminal or shell without specifying its full path.

## Behavior

If no config is specified, `backscribe` searches upward from the current directory for:

- backscribe.{yaml,yml,json}
- .backscribe.{yaml,yml,json}

It also checks `$HOME` and `$HOME/.config` directories.

## Configuration Format

Your configuration defines groups of file operations. Each group contains archives that describe what to back up or copy.

Archive Types:

- 7z: Create a .7z archive
- tar: Create a .tar archive
- cp: Copy files directly

## Example Configs

✅ YAML

```yaml
groups:
  - name: home
    archives:
      - type: 7z
        name: dotfiles
        files:
          - ~/.*
        exclude:
          - .cache
```

✅ JSON

```json
{
  "groups": [
    {
      "name": "home",
      "archives": [
        {
          "type": "7z",
          "name": "dotfiles",
          "files": ["~/.*"],
          "exclude": [".cache"]
        }
      ]
    }
  ]
}
```

## Environment Variable Expansion

You can use environment variables in paths, like:

```yaml
groups:
  - name: roaming
    archives:
      - type: 7z
        name: Roaming
        files:
          - ${APPDATA}/Code
```

## Visual Studio Code

See `.vscode/settings.json` for enabling schema validation, auto-completion, and hover tooltips with documentation for *backscribe* files.

## Building from Source

```sh
git clone https://github.com/akornatskyy/backscribe.git
cd backscribe
go build -ldflags="-s -w"
```
