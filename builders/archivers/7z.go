package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func Build7z(archive domain.Archive, group domain.Group) string {
	var src []string

	options := []string{"-t7z", "-bso0"}
	if archive.Method != nil && archive.Method.Level != nil {
		options = append(options, fmt.Sprintf("-mx%d", *archive.Method.Level))
	}

	src = append(src, fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.7z" ]; then
    log '\e[32m▶\e[0m %s => %s'
    7z a %s "${DEST_DIR}/%s.7z" \
      `,
		archive.Name,
		group.Name, archive.Name,
		strings.Join(options, " "), archive.Name,
	))

	files := []string{}
	for _, f := range archive.Files {
		files = append(files, helpers.Quote(f))
	}
	for _, x := range archive.Exclude {
		files = append(files, "-x!"+helpers.Quote(x))
	}
	for _, x := range archive.Rexclude {
		files = append(files, "-xr!"+helpers.Quote(x))
	}

	src = append(src, strings.Join(files, " \\\n      "))
	src = append(src, `
  else
    log '\e[33m↷\e[0m `+group.Name+` => `+archive.Name+`'
  fi
`)

	return strings.Join(src, "")
}
