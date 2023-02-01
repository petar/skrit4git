package cmd

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

const (
	AgentName           = "twitter4git"
	AgentVarPath        = "." + AgentName
	AgentConfigFilebase = "." + AgentName + ".json"
	AgentTempPath       = AgentName
)

type Setup struct {
	Repo git.Address
}

type Config struct {
	RepoURL    git.URL    `json:"repo_url"`
	RepoBranch git.Branch `json:"repo_branch"`
	Auth       AuthConfig `json:"auth"`
	VarDir     string     `json:"var_dir"`
}

type AuthConfig struct {
	SSHPrivateKeysFile *string       `json:"ssh_private_keys_file"`
	AccessToken        *string       `json:"access_token"`
	UserPassword       *UserPassword `json:"user_password"`
}

type UserPassword struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (cfg Config) Setup(ctx context.Context) Setup {

	git.SetAuthor("twitter4git agent", "no-reply@twitter4git.xyz")

	switch {
	case cfg.Auth.SSHPrivateKeysFile != nil:
		git.SetAuth(ctx, cfg.RepoURL, git.MakeSSHFileAuth(ctx, "git", *cfg.Auth.SSHPrivateKeysFile))
	case cfg.Auth.AccessToken != nil:
		git.SetAuth(ctx, cfg.RepoURL, git.MakeTokenAuth(ctx, *cfg.Auth.AccessToken))
	case cfg.Auth.UserPassword != nil:
		git.SetAuth(ctx, cfg.RepoURL, git.MakePasswordAuth(ctx, cfg.Auth.UserPassword.User, cfg.Auth.UserPassword.Password))
	}

	return Setup{
		Repo: git.Address{Repo: cfg.RepoURL, Branch: cfg.RepoBranch},
	}
}
