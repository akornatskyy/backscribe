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
    log '\e[32m↻\e[0m `+group.Name+` => `+archive.Name+`'`)
	src = append(src, fmt.Sprintf(`
    rm -f "${DEST_DIR}/%s.up.tar.tmp"`, archive.Name))
	src = append(src, fmt.Sprintf(`
    %s u -ttar -bso0 \
      "${DEST_DIR}/%s.tar" -u- \
      -up0q3x2y2z0!"${DEST_DIR}/%s.up.tar.tmp" \
      `,
		cmd,
		archive.Name,
		archive.Name))
	src = append(src, strings.Join(files, " \\\n      "))
	src = append(src, fmt.Sprintf(`
    if [ "$(stat -c%%s "${DEST_DIR}/%s.up.tar.tmp")" -ne 1024 ]; then
      mv "${DEST_DIR}/%s.up.tar.tmp" \
        "${DEST_DIR}/%s.up.tar"
    else
      echo -n -e "\e[1A\e[K"
      log '\e[33m↷\e[0m %s => %s'
      rm "${DEST_DIR}/%s.up.tar.tmp"
    fi
  fi
`, archive.Name,
		archive.Name,
		archive.Name,
		group.Name, archive.Name,
		archive.Name))

	return strings.Join(src, "")
}
