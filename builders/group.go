package builders

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/domain"
)

func BuildGroup(group domain.Group) string {
	src := []string{fmt.Sprintf("\nbackup_%s() {", group.Name)}
	for _, archive := range group.Archives {
		chunk := BuildArchive(archive, group)
		if chunk != "" {
			src = append(src, chunk)
		}
	}
	if len(src) == 1 {
		src = append(src, "\n  :\n")
	}
	src = append(src, "}")
	return strings.Join(src, "")
}
