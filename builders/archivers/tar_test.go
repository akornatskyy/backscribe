package archivers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akornatskyy/backscribe/builders/archivers"
	"github.com/akornatskyy/backscribe/domain"
)

func TestBuildTar_Empty(t *testing.T) {
	got := archivers.BuildTar(domain.Archive{}, domain.Group{})
	if got != "" {
		t.Errorf("Expected empty; got: %s", got)
	}
}

func TestBuildTar_GoldenCases(t *testing.T) {
	tests := []struct {
		name    string
		archive domain.Archive
		group   domain.Group
	}{
		{
			name: "files-only",
			archive: domain.Archive{
				Name:  "files-only",
				Files: []string{"file 1.txt", "file2.txt"},
			},
			group: domain.Group{Name: "group1"},
		},
		{
			name: "full",
			archive: domain.Archive{
				Name:     "full",
				Cwd:      "/project path",
				Files:    []string{"file 1.txt", "file2.txt"},
				Exclude:  []string{"temp/", "*.tmp"},
				Rexclude: []string{"logs/**"},
				Include:  []string{"*.go"},
				Rinclude: []string{"**/*.md"},
			},
			group: domain.Group{Name: "group2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := archivers.BuildTar(tt.archive, tt.group)
			goldenPath := filepath.Join("testdata", "tar-"+tt.name+".golden")

			wantBytes, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("Failed to read golden file %s: %v", goldenPath, err)
			}

			want := string(wantBytes)
			if got != want {
				t.Errorf("Mismatch with golden file %s\n--- Got:\n%s\n--- Want:\n%s",
					goldenPath, got, want)
			}
		})
	}
}
