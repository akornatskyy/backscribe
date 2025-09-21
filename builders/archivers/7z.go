package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func Build7z(archive domain.Archive, group domain.Group) string {
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

	options := []string{"-t7z", "-bso0"}
	if archive.Method != nil && archive.Method.Level != nil {
		options = append(options, fmt.Sprintf("-mx%d", *archive.Method.Level))
	}

	return fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.7z" ]; then
    log '\e[32m▶\e[0m %s => %s'
    %s a %s \
      "${DEST_DIR}/%s.7z" \
      %s
  else
    log '\e[32m↻\e[0m %s => %s'
    rm -f "${DEST_DIR}/%s.up.7z.tmp"
    %s u %s \
      "${DEST_DIR}/%s.7z" -u- \
      -up0q3x2y2z0!"${DEST_DIR}/%s.up.7z.tmp" \
      %s
    if [ "$(stat -c%%s "${DEST_DIR}/%s.up.7z.tmp")" -ne 32 ]; then
      mv "${DEST_DIR}/%s.up.7z.tmp" \
        "${DEST_DIR}/%s.up.7z"
    else
      printf '\033[1A\033[K'
      log '\e[33m↷\e[0m %s => %s'
      rm "${DEST_DIR}/%s.up.7z.tmp"
    fi
  fi
`,
		archive.Name,
		group.Name, archive.Name,
		cmd, strings.Join(options, " "),
		archive.Name,
		strings.Join(args, " \\\n      "),
		group.Name, archive.Name,
		archive.Name,
		cmd, strings.Join(options, " "),
		archive.Name,
		archive.Name, strings.Join(args, " \\\n      "),
		archive.Name, archive.Name,
		archive.Name,
		group.Name, archive.Name,
		archive.Name,
	)
}

type SevenZipBuilder struct{}

func (b SevenZipBuilder) Build(
	archive domain.Archive, group domain.Group) string {
	return Build7z(archive, group)
}

func init() {
	RegisterBuilder("7z", SevenZipBuilder{})
}
