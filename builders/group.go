package builders

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/domain"
)

func BuildGroup(group domain.Group) string {
	src := []string{fmt.Sprintf("\nbackup_%s() {", group.Name)}
	for _, archive := range group.Archives {
		src = append(src, BuildArchive(archive, group))
	}
	src = append(src, "}")
	return strings.Join(src, "")
}
