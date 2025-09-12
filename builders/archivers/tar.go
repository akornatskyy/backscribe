package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildTar(archive domain.Archive, group domain.Group) string {
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
	for _, i := range archive.Include {
		files = append(files, "-i!"+helpers.Quote(i))
	}
	for _, i := range archive.Rinclude {
		files = append(files, "-ir!"+helpers.Quote(i))
	}

	if len(files) == 0 {
		return ""
	}

	cmd := "7z"
	if archive.Cwd != "" {
		cmd = "cd " + helpers.Quote(archive.Cwd) + " && " + cmd
	}

	src := []string{fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.tar" ]; then
    log '\e[32m▶\e[0m %s => %s'
    %s a -ttar -bso0 \
      "${DEST_DIR}/%s.tar" \
      `,
		archive.Name,
		group.Name, archive.Name,
		cmd,
		archive.Name,
	)}

	src = append(src, strings.Join(files, " \\\n      "))
	src = append(src, `
  else
    log '\e[33m↷\e[0m `+group.Name+` => `+archive.Name+`'
  fi
`)

	return strings.Join(src, "")
}
