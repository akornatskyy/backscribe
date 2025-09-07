package helpers

import (
	"strings"
)

func Quote(name string) string {
	if (strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`)) ||
		(strings.HasPrefix(name, `'`) && strings.HasSuffix(name, `'`)) {
		return name
	}
	if strings.Contains(name, `\ `) {
		return name
	}
	if !strings.Contains(name, " ") {
		return name
	}
	if strings.HasPrefix(name, "~") {
		return strings.ReplaceAll(name, ` `, `\ `)
	}
	return `"` + strings.ReplaceAll(name, `"`, `\"`) + `"`
}
