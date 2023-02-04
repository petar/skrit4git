package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"time"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

type Handle string // user's public handle, which is a host/path pair

func ParseHandle(repo git.URL) (Handle, error) {
	u, err := url.Parse(string(repo))
	if err != nil {
		return "", err
	}
	if u.Port() == "" {
		return Handle(u.Host + "/" + u.Path), nil
	} else {
		return Handle(u.Host + ":" + u.Port() + "/" + u.Path), nil
	}
}

type Home struct {
	Handle  Handle
	Public  git.Address
	Private git.Address
}

var RootNS = ns.NS{}

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, ProtocolName+": "+msg)
}

const (
	ProtocolName           = "twitter4git"
	ProtocolVersion        = "0.0.1"
	PostDir                = "post"
	PostFilenameTimeFormat = "20060102-150405"
)

type LocalID string // YYYYMMDD-HHMMSS-SHA256CONTENT

func (x LocalID) String() string {
	return string(x)
}

func PostNS(t time.Time, content string) (ns.NS, LocalID) {
	localID := PostFilebase(t, content)
	return RootNS.Join(ns.NS{PostDir, localID.String()}), localID
}

// PostFilebase returns a filename of the form YYYYMMDD-HHMMSS-SHA256CONTENT
func PostFilebase(t time.Time, content string) LocalID {
	return LocalID(t.UTC().Format(PostFilenameTimeFormat) + "-" + ContentHash(content))
}

func ContentHash(content string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(content)); err != nil {
		panic(err)
	}
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

type PostMeta struct {
	By Handle `json:"by"`
}
