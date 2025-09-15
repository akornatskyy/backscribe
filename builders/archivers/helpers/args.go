package helpers

import "github.com/akornatskyy/backscribe/domain"

func Build7zArgsFromArchive(a domain.Archive) []string {
	var args []string
	appendArgs := func(prefix string, items []string) {
		for _, v := range items {
			if v != "" {
				args = append(args, prefix+Quote(v))
			}
		}
	}
	appendArgs("", a.Files)
	appendArgs("-x!", a.Exclude)
	appendArgs("-xr!", a.Rexclude)
	appendArgs("-i!", a.Include)
	appendArgs("-ir!", a.Rinclude)
	return args
}
