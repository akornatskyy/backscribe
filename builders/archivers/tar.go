package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildTar(archive domain.Archive, group domain.Group) string {
	src := []string{fmt.Sprintf(`
  if [ ! -e "${DEST_DIR}/%s.tar" ]; then
    log '\e[32m▶\e[0m %s => %s'
    7z a -ttar -bso0 "${DEST_DIR}/%s.tar" \
      `,
		archive.Name, group.Name, archive.Name, archive.Name,
	)}

	filesAndExcludes := []string{}
	for _, f := range archive.Files {
		filesAndExcludes = append(filesAndExcludes, helpers.Quote(f))
	}
	if archive.Exclude != nil {
		for _, x := range archive.Exclude {
			filesAndExcludes = append(filesAndExcludes, "-x!"+helpers.Quote(x))
		}
	}
	if archive.Rexclude != nil {
		for _, x := range archive.Rexclude {
			filesAndExcludes = append(filesAndExcludes, "-xr!"+helpers.Quote(x))
		}
	}

	src = append(src, strings.Join(filesAndExcludes, " \\\n      "))
	src = append(src, `
  else
    log '\e[33m↷\e[0m `+group.Name+` => `+archive.Name+`'
  fi
`)

	return strings.Join(src, "")
}
