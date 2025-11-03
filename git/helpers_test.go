package git

import (
	"testing"

	"github.com/ejoffe/spr/config"
)

func TestBranchNameRegex(t *testing.T) {
	cfg := config.DefaultConfig()
	tests := []struct {
		prefix string
		input  string
		branch string
		commit string
	}{
		{prefix: "", input: "spr/b1/deadbeef", branch: "b1", commit: "deadbeef"},
		{prefix: "team-a/", input: "team-a/spr/b1/de4dbeef", branch: "q-w", commit: "de4dbeef"},
	}

	for _, tc := range tests {
		cfg.Repo.PrBranchPrefix = tc.prefix
		commitID := CommitIDFromBranchName(cfg, tc.input)
		if commitID == nil {
			t.Fatalf("expected commit ID to be extracted from %q, but got nil", tc.input)
		}
		if *commitID != tc.commit {
			t.Fatalf("expected: '%v', actual: '%v'", tc.commit, *commitID)
		}
	}
}

func TestBranchNameRegexNegative(t *testing.T) {
	cfg := config.DefaultConfig()
	tests := []struct {
		prefix string
		input  string
	}{
		{prefix: "", input: "feature/b1/deadbeef"},
		{prefix: "", input: "team-b/spr/b1/deadbeef"},
		{prefix: "team-a/", input: "team-a/feature/b1/deadbeef"},
		{prefix: "team-a/", input: "spr/b1/deadbeef"},
	}

	for _, tc := range tests {
		cfg.Repo.PrBranchPrefix = tc.prefix
		commitID := CommitIDFromBranchName(cfg, tc.input)
		if commitID != nil {
			t.Fatalf("expected no commit ID to be extracted from %q, but got %q", tc.input, *commitID)
		}
	}
}
