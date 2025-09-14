package archivers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akornatskyy/backscribe/domain"
	"github.com/akornatskyy/backscribe/builders/archivers"
)

func TestBuildCopy_Empty(t *testing.T) {
	got := archivers.BuildCopy(domain.Archive{}, domain.Group{})
	if got != "" {
		t.Errorf("Expected empty; got: %s", got)
	}
}

func TestBuildCopy_GoldenCases(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := archivers.BuildCopy(tt.archive, tt.group)
			goldenPath := filepath.Join("testdata", "cp-"+tt.name+".golden")

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
