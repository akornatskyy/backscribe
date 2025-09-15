package helpers_test

import (
	"reflect"
	"testing"

	"github.com/akornatskyy/backscribe/builders/archivers/helpers"
	"github.com/akornatskyy/backscribe/domain"
)

func TestBuild7zArgsFromArchive(t *testing.T) {
	tests := []struct {
		name     string
		archive  domain.Archive
		expected []string
	}{
		{
			name:     "all fields empty",
			archive:  domain.Archive{},
			expected: nil,
		},
		{
			name: "only files",
			archive: domain.Archive{
				Files: []string{"file1.txt", "file2.txt"},
			},
			expected: []string{
				"file1.txt",
				"file2.txt",
			},
		},
		{
			name: "exclude and reexclude",
			archive: domain.Archive{
				Exclude:  []string{"temp/*"},
				Rexclude: []string{"*.log"},
			},
			expected: []string{
				"-x!temp/*",
				"-xr!*.log",
			},
		},
		{
			name: "include and rinclude",
			archive: domain.Archive{
				Include:  []string{"*.go"},
				Rinclude: []string{"src/*"},
			},
			expected: []string{
				"-i!*.go",
				"-ir!src/*",
			},
		},
		{
			name: "mixed with empty strings",
			archive: domain.Archive{
				Files:    []string{"main.go", ""},
				Exclude:  []string{"", "node_modules"},
				Rexclude: []string{"*.tmp", ""},
				Include:  []string{"*.md"},
				Rinclude: []string{""},
			},
			expected: []string{
				"main.go",
				"-x!node_modules",
				"-xr!*.tmp",
				"-i!*.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.Build7zArgsFromArchive(tt.archive)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
