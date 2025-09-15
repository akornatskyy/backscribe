package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildTar(archive domain.Archive, group domain.Group) string {
	args := helpers.Build7zArgsFromArchive(archive)
	if len(args) == 0 {
		return ""
	}

	var cmd string
	if archive.Cwd != "" {
		cmd = fmt.Sprintf("cd %s && 7z", helpers.Quote(archive.Cwd))
	} else {
		cmd = "7z"
	}

	return fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.tar" ]; then
    log '\e[32m▶\e[0m %s => %s'
    %s a -ttar -bso0 \
      "${DEST_DIR}/%s.tar" \
      %s
  else
    log '\e[32m↻\e[0m %s => %s'
    rm -f "${DEST_DIR}/%s.up.tar.tmp"
    %s u -ttar -bso0 \
      "${DEST_DIR}/%s.tar" -u- \
      -up0q3x2y2z0!"${DEST_DIR}/%s.up.tar.tmp" \
      %s
    if [ "$(stat -c%%s "${DEST_DIR}/%s.up.tar.tmp")" -ne 1024 ]; then
      mv "${DEST_DIR}/%s.up.tar.tmp" \
        "${DEST_DIR}/%s.up.tar"
    else
      echo -n -e "\e[1A\e[K"
      log '\e[33m↷\e[0m %s => %s'
      rm "${DEST_DIR}/%s.up.tar.tmp"
    fi
  fi
`,
		archive.Name,
		group.Name, archive.Name,
		cmd,
		archive.Name,
		strings.Join(args, " \\\n      "),
		group.Name, archive.Name,
		archive.Name,
		cmd,
		archive.Name,
		archive.Name, strings.Join(args, " \\\n      "),
		archive.Name,
		archive.Name,
		archive.Name,
		group.Name, archive.Name,
		archive.Name,
	)
}

type TarBuilder struct{}

func (b TarBuilder) Build(
	archive domain.Archive, group domain.Group) string {
	return BuildTar(archive, group)
}

func init() {
	RegisterBuilder("tar", TarBuilder{})
}
