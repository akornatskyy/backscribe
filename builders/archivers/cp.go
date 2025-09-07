package archivers

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildCopy(archive domain.Archive, group domain.Group) string {
	var src []string

	target := ""
	if archive.Dst != "" {
		target = "/" + archive.Dst
		src = append(src, fmt.Sprintf("\n  mkdir -p \"${DEST_DIR}%s\"", target))
	}

	files := make([]string, len(archive.Src))
	for i, f := range archive.Src {
		files[i] = helpers.Quote(f)
	}

	src = append(src, fmt.Sprintf(`
  log '\e[32m▶\e[0m %s => %s'
  cp -a -n %s "${DEST_DIR}%s"
`, group.Name, archive.Name, strings.Join(archive.Src, " "), target))

	return strings.Join(src, "")
}
