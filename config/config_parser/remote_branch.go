package config_parser

import (
	"regexp"
	"strings"

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

/*
Parses the local branch and the remote ref from the output of `git

	status -b --porcelain -u no`.  The second capture contains the
	entire remote ref (remote/branch...).
*/
var _remoteBranchRegex = regexp.MustCompile(`^## ([A-Za-z0-9_./-]+)\.\.\.([^\s]+)`)

func (s *remoteBranch) Load(cfg interface{}) {
	var output string
	err := s.gitcmd.Git("status -b --porcelain -u no", &output)
	check(err)

	matches := _remoteBranchRegex.FindStringSubmatch(output)
	if matches == nil {
		return
	}

	repoCfg := cfg.(*config.RepoConfig)

	/* matches[2] is e.g. "myremote/myuser/mybranch".  Split at the
	   first '/' to get the remote name and the remote branch.  */
	remoteRef := matches[2]
	parts := strings.SplitN(remoteRef, "/", 2)

	repoCfg.GitHubRemote = parts[0]
	if len(parts) > 1 {
		repoCfg.GitHubBranch = parts[1]
	} else {
		repoCfg.GitHubBranch = ""
	}
}
