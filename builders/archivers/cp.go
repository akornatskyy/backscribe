package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildCopy(archive domain.Archive, group domain.Group) string {
	files := make([]string, len(archive.Files))
	for i, f := range archive.Files {
		if f != "" {
			files[i] = helpers.Quote(f)
		}
	}
	if len(files) == 0 {
		return ""
	}

	var cmd string
	if archive.Cwd != "" {
		cmd = fmt.Sprintf("cd %s && cp", helpers.Quote(archive.Cwd))
	} else {
		cmd = "cp"
	}

	return fmt.Sprintf(`
  if [ ! -d "${DEST_DIR}/%s" ]; then
    log '\e[32m▶\e[0m %s => %s'
    mkdir -p "${DEST_DIR}/%s"
  else
    log '\e[32m↻\e[0m %s => %s'
  fi
  %s -a -u -t "${DEST_DIR}/%s" \
    %s
`,
		archive.Name,
		group.Name, archive.Name,
		archive.Name,
		group.Name, archive.Name,
		cmd, archive.Name,
		strings.Join(files, " \\\n    "))
}

type CopyBuilder struct{}

func (b CopyBuilder) Build(
	archive domain.Archive, group domain.Group) string {
	return BuildCopy(archive, group)
}

func init() {
	RegisterBuilder("cp", CopyBuilder{})
}
