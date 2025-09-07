package builders

import (
	"github.com/akornatskyy/backscribe/builders/archivers"
	"github.com/akornatskyy/backscribe/domain"
)

func BuildArchive(archive domain.Archive, group domain.Group) string {
	switch archive.Type {
	case "7z":
		return archivers.Build7z(archive, group)
	case "tar":
		return archivers.BuildTar(archive, group)
	case "cp":
		return archivers.BuildCopy(archive, group)
	default:
		return ""
	}
}
