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
		files[i] = helpers.Quote(f)
	}

	return fmt.Sprintf(`
	if [ ! -d "${DEST_DIR}/%s" ]; then
  	log '\e[32m▶\e[0m [cp] %s => %s'
		mkdir -p "${DEST_DIR}/%s"
  	cp -a -n -t "${DEST_DIR}/%s" \
			%s
	else
    log '\e[33m↷\e[0m %s => %s'
  fi
`, archive.Name,
		group.Name, archive.Name,
		archive.Name,
		archive.Name,
		strings.Join(files, " \\\n      "),
		group.Name, archive.Name)
}
