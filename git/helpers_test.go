package git

import (
	"testing"

	"github.com/ejoffe/spr/config"
)

func TestBranchNameRegex(t *testing.T) {
	cfg := config.DefaultConfig()
	tests := []struct {
		input  string
		branch string
		commit string
	}{
		{input: "spr/b1/deadbeef", branch: "b1", commit: "deadbeef"},
	}

	for _, tc := range tests {
		commitID := CommitIDFromBranchName(cfg, tc.input)
		if commitID == nil {
			t.Fatalf("expected commit ID to be extracted from %q, but got nil", tc.input)
		}
		if *commitID != tc.commit {
			t.Fatalf("expected: '%v', actual: '%v'", tc.commit, *commitID)
		}
	}
}
