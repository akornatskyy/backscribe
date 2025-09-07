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
		archive.Name, group.Name, archive.Name,
		strings.Join(options, " "), archive.Name,
	))

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
	if archive.Copy != nil {
		for _, x := range archive.Copy {
			filesAndExcludes = append(filesAndExcludes, "-xr!"+helpers.Quote(x))
		}
	}

	src = append(src, strings.Join(filesAndExcludes, " \\\n      "))

	if len(archive.Copy) > 0 {
		options = helpers.FilterOut(options, func(o string) bool {
			return !strings.HasPrefix(o, "-mx")
		})
		options = append(options, "-mx0")

		src = append(src, fmt.Sprintf(`

    7z u %s "${DEST_DIR}/%s.7z" \
      `, strings.Join(options, " "), archive.Name))

		filesAndExcludes = []string{}
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
		if archive.Copy != nil {
			for _, x := range archive.Copy {
				filesAndExcludes = append(filesAndExcludes, "-ir!"+helpers.Quote(x))
			}
		}

		src = append(src, strings.Join(filesAndExcludes, " \\\n      "))
	}

	src = append(src, `
  else
    log '\e[33m↷\e[0m `+group.Name+` => `+archive.Name+`'
  fi
`)

	return strings.Join(src, "")
}
