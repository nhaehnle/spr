package config_parser

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ejoffe/rake"
	"github.com/ejoffe/spr/config"
	"github.com/ejoffe/spr/git"
)

func ParseConfig(gitcmd git.GitInterface) *config.Config {
	cfg := config.EmptyConfig()

	repoCommonConfigPath, repoRemoteConfigPath := getRepoConfigFilePath(gitcmd)

	fmt.Println("common", repoCommonConfigPath)
	fmt.Println("remote", repoRemoteConfigPath)

	rake.LoadSources(cfg.Repo,
		rake.DefaultSource(),
		NewGitHubRemoteSource(cfg, gitcmd),
		NewRemoteBranchSource(gitcmd),
	)

	_, err := os.Stat(repoCommonConfigPath)
	if err == nil {
		rake.LoadSources(cfg.Repo,
			rake.YamlFileSource(repoCommonConfigPath),
		)
	}

	if repoRemoteConfigPath != "" {
		_, err = os.Stat(repoRemoteConfigPath)
		if err == nil {
			rake.LoadSources(cfg.Repo,
				rake.YamlFileSource(repoRemoteConfigPath),
			)
		}
	}

	if cfg.Repo.GitHubHost == "" {
		fmt.Println("unable to auto configure repository host - must be set manually in .spr.yml or .{remote}.spr.yml")
		os.Exit(2)
	}
	if cfg.Repo.GitHubRepoOwner == "" {
		fmt.Println("unable to auto configure repository owner - must be set manually in .spr.yml or .{remote}.spr.yml")
		os.Exit(3)
	}

	if cfg.Repo.GitHubRepoName == "" {
		fmt.Println("unable to auto configure repository name - must be set manually in .spr.yml or .{remote}.spr.yml")
		os.Exit(4)
	}

	rake.LoadSources(cfg.User,
		rake.DefaultSource(),
		rake.YamlFileSource(UserConfigFilePath()),
	)

	rake.LoadSources(cfg.State,
		rake.DefaultSource(),
		rake.YamlFileSource(InternalConfigFilePath()),
	)

	cfg.State.RunCount = cfg.State.RunCount + 1

	rake.LoadSources(cfg.State,
		rake.YamlFileWriter(InternalConfigFilePath()))

	return cfg
}

func CheckConfig(cfg *config.Config) error {
	return nil
}

func getRepoConfigFilePath(gitcmd git.GitInterface) (commonConfig string, remoteConfig string) {
	rootdir := gitcmd.RootDir()
	commonConfig = filepath.Clean(path.Join(rootdir, ".spr.yml"))

	remote, _, err := git.GetTrackedUpstream(gitcmd)
	if err == nil {
		remoteConfig = filepath.Clean(path.Join(rootdir, fmt.Sprintf(".%s.spr.yml", remote)))
	} else {
		remoteConfig = ""
	}

	return commonConfig, remoteConfig
}

func UserConfigFilePath() string {
	rootdir, err := os.UserHomeDir()
	check(err)
	filepath := filepath.Clean(path.Join(rootdir, ".spr.yml"))
	return filepath
}

func InternalConfigFilePath() string {
	rootdir, err := os.UserHomeDir()
	check(err)
	filepath := filepath.Clean(path.Join(rootdir, ".spr.state"))
	return filepath
}
