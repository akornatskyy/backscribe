package builders

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/akornatskyy/backscribe/domain"
)

func BuildScript(config *domain.Config, configFile string) string {
	src := []string{`#!/bin/sh
set -o errexit
# set -o xtrace
`}
	src = append(src, fmt.Sprintf(`
CONFIG_FILE=%s
H=$(cd ~ ; (pwd -W 2>/dev/null || pwd) | sed 's/:\//:\/\//')
DEST_DIR="${H}/backups/$(date '+%%Y-%%m-%%d')"
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
		src = append(src, BuildGroup(group))
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
			src = append(src, fmt.Sprintf("\nlog '\\e[36m⚬\\e[0m %s'", group.Name))
		}
		src = append(src, fmt.Sprintf("\n%sbackup_%s", prefix, group.Name))
	}

	src = append(src, "\n\nlog '\\e[32m✓\\e[0m all done'")
	return strings.Join(src, "")
}
