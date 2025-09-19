package builders

import (
	"fmt"
	"strings"

	"github.com/akornatskyy/backscribe/domain"
)

func BuildGroup(b *strings.Builder, group domain.Group) {
	fmt.Fprintf(b, "\nbackup_%s() {", sanitizeName(group.Name))
	added := false
	for _, archive := range group.Archives {
		if chunk := BuildArchive(archive, group); chunk != "" {
			b.WriteString(chunk)
			added = true
		}
	}

	if !added {
		b.WriteString("\n  :\n")
	}
	b.WriteString("}\n")
}

func sanitizeName(name string) string {
	var sb strings.Builder
	sb.Grow(len(name))

	for _, r := range name {
		switch {
		case r == '_',
			r >= '0' && r <= '9',
			r >= 'A' && r <= 'Z',
			r >= 'a' && r <= 'z':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	return sb.String()
}
