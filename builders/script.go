package builders

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/akornatskyy/backscribe/domain"
)

func BuildScript(config *domain.Config, configFile string) string {
	var sb strings.Builder
	sb.Grow(estimateScriptSize(config))

	fmt.Fprintf(&sb, `#!/bin/sh
set -o errexit
# set -o xtrace

CONFIG_FILE=%s
H=$(cd ~ ; (pwd -W 2>/dev/null || pwd) | sed 's/:\//:\/\//')
: "${BACKSCRIBE_BACKUPS_DIR:="${H}/backups"}"
DEST_DIR=${BACKSCRIBE_BACKUPS_DIR}
START=$(date +%%s)

log() {
  NOW=$(date +%%s)
  ELAPSED=$((NOW - START))
  MIN=$((ELAPSED / 60))
  SEC=$((ELAPSED %% 60))
  if [ "${MIN}" -lt 10 ]; then MIN="0${MIN}"; fi
  if [ "${SEC}" -lt 10 ]; then SEC="0${SEC}"; fi
  echo -e "\033[90m${MIN}:${SEC}\033[0m" $1
}
`, strconv.Quote(configFile))

	for _, group := range config.Groups {
		BuildGroup(&sb, group)
	}

	sb.WriteString(`
echo "CONFIG_FILE=${CONFIG_FILE}"
echo "DEST_DIR=${DEST_DIR}"
mkdir -p "${DEST_DIR}"

`)

	for _, group := range config.Groups {
		if group.Skip {
			fmt.Fprintf(&sb,
				"log '\\e[36m⚬\\e[0m %s'\n# backup_%s\n",
				group.Name, group.Name)
		} else {
			fmt.Fprintf(&sb, "backup_%s\n", group.Name)
		}
	}

	sb.WriteString("\nlog '\\e[32m✓\\e[0m all done'")
	return sb.String()
}

func estimateScriptSize(config *domain.Config) int {
	size := 700
	for _, group := range config.Groups {
		size += len(group.Archives) * 1000
	}
	return size
}
