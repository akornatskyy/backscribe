# backscribe

[![tests](https://github.com/akornatskyy/backscribe/actions/workflows/tests.yml/badge.svg)](https://github.com/akornatskyy/backscribe/actions/workflows/tests.yml)

*Backscribe* is a flexible command generator for archiving, copying, and backing up files, driven by structured configuration definitions. It outputs a series of shell commands that you can review or pipe directly into your shell (e.g., via `sh` or `bash`).

## Usage

```text
Usage:
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
  backscribe | sh  # Auto-searches for config
```

## Prerequisites

The 7z ([7-Zip](https://sourceforge.net/projects/sevenzip/files/7-Zip/)) command-line tool must be installed and discoverable in your system’s PATH. You should be able to run 7z from any terminal or shell without specifying its full path.

## Behavior

If no config is specified, `backscribe` searches upward from the current directory for:

- backscribe.{yaml,yml,json}
- .backscribe.{yaml,yml,json}

It also checks `$HOME` and `$HOME/.config` directories.

### Environment Variables

#### `BACKSCRIBE_BACKUPS_DIR`

Default: `~/backups`

Specifies the destination directory for archive files. If set, this variable overrides the default location.

```sh
export BACKSCRIBE_BACKUPS_DIR="/mnt/storage/backups"
```

In this example, all archives will be written to `/mnt/storage/backups` instead of `~/backups`.

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

To enable schema validation, auto-completion, and hover tooltips with documentation for *backscribe* files:

```json
{
  "json.schemas": [
    {
      "fileMatch": ["backscribe.json", ".backscribe.json"],
      "url": "https://raw.githubusercontent.com/akornatskyy/backscribe/refs/heads/main/schema.json"
    }
  ],
  "yaml.schemas": {
    "https://raw.githubusercontent.com/akornatskyy/backscribe/refs/heads/main/schema.json": [
      "backscribe.yaml",
      "backscribe.yml",
      ".backscribe.yaml",
      ".backscribe.yml"
    ]
  }
}
```

## Building from Source

```sh
git clone https://github.com/akornatskyy/backscribe.git
cd backscribe
go build
```

## Troubleshooting

### Duplicate filename on disk

If you see an error like this:

```txt
ERROR:
Duplicate filename on disk:
.bash_history
.bash_history
```

it means the `7z` command tried to add the same file (in this case, `.bash_history`) more than once when creating the archive.

This usually occurs because of how shell globbing patterns expand. For example, using `~/.*` will match:

- `.` (the current directory)
- `..` (the parent directory)
- all hidden files (dotfiles)

As a result, some files may be included multiple times or in overlapping ways, which causes the duplication error.

Instead of `~/.*`, use a safer globbing pattern that excludes `.` and `..`. For example:

```yaml
groups:
  - name: home
    archives:
      - type: 7z
        name: dotfiles
        files:
          - ~/.[!.]*   # matches hidden files, but not '.' or '..'
        exclude:
          - .cache     # exclude large/unnecessary directories
```

This pattern ensures only actual hidden files are included, avoiding duplicates.
