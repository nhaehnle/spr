package config_parser

import (
	"github.com/ejoffe/spr/config"
	"github.com/ejoffe/spr/git"
)

type remoteBranch struct {
	gitcmd git.GitInterface
}

func NewRemoteBranchSource(gitcmd git.GitInterface) *remoteBranch {
	return &remoteBranch{
		gitcmd: gitcmd,
	}
}

func (s *remoteBranch) Load(cfg interface{}) {
	remote, branch, err := git.GetTrackedUpstream(s.gitcmd)
	if err != nil {
		return
	}

	repoCfg := cfg.(*config.RepoConfig)
	repoCfg.GitHubRemote = remote
	repoCfg.GitHubBranch = branch
}
