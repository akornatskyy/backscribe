package archivers

import "github.com/akornatskyy/backscribe/domain"

type ArchiveBuilder interface {
	Build(archive domain.Archive, group domain.Group) string
}

var registry = map[string]ArchiveBuilder{}

func RegisterBuilder(name string, builder ArchiveBuilder) {
	registry[name] = builder
}

func GetBuilder(name string) ArchiveBuilder {
	return registry[name]
}
