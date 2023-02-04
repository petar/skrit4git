package cmd

import (
	"context"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/petar/twitter4git/proto"
)

const (
	AgentName           = "twitter4git"
	AgentVarPath        = "." + AgentName
	AgentConfigFilebase = "." + AgentName + ".json"
	AgentTempPath       = AgentName
)

type Setup struct {
	Home proto.Home
}

type Config struct {
	PublicHomeURL  git.URL    `json:"public_home_url"`  // public URL of home repo
	PrivateHomeURL git.URL    `json:"private_home_url"` // private URL of home repo
	Auth           AuthConfig `json:"auth"`
	VarDir         string     `json:"var_dir"`
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

	git.SetAuthor(proto.ProtocolName+" agent", "no-reply@"+proto.ProtocolName+".xyz")

	switch {
	case cfg.Auth.SSHPrivateKeysFile != nil:
		git.SetAuth(ctx, cfg.PrivateHomeURL, git.MakeSSHFileAuth(ctx, "git", *cfg.Auth.SSHPrivateKeysFile))
	case cfg.Auth.AccessToken != nil:
		git.SetAuth(ctx, cfg.PrivateHomeURL, git.MakeTokenAuth(ctx, *cfg.Auth.AccessToken))
	case cfg.Auth.UserPassword != nil:
		git.SetAuth(ctx, cfg.PrivateHomeURL, git.MakePasswordAuth(ctx, cfg.Auth.UserPassword.User, cfg.Auth.UserPassword.Password))
	}

	handle, err := proto.ParseHandle(cfg.PublicHomeURL)
	must.NoError(ctx, err)
	return Setup{
		Home: proto.Home{
			Handle:  handle,
			Public:  git.Address{Repo: cfg.PublicHomeURL, Branch: ""},  // the app will create sub-branches
			Private: git.Address{Repo: cfg.PrivateHomeURL, Branch: ""}, // the app will create sub-branches
		},
	}
}
