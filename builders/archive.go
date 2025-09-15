package builders

import (
	"github.com/akornatskyy/backscribe/builders/archivers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildArchive(archive domain.Archive, group domain.Group) string {
	builder := archivers.GetBuilder(archive.Type)
	if builder == nil {
		return ""
	}
	return builder.Build(archive, group)
}
