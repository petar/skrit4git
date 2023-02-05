package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

type Handle string // user's public handle, which is a host/path pair

func (h Handle) String() string {
	return string(h)
}

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
	Public  git.Address // public read-only URL to home repo + branch prefix
	Private git.Address // private read/write URL to home repo + branch prefix
}

func (h Home) PublicSend() git.Address {
	return h.Public.Sub(SendBranchSuffix)
}

func (h Home) PublicReceive() git.Address {
	return h.Public.Sub(ReceiveBranchSuffix)
}

func (h Home) PrivateSend() git.Address {
	return h.Private.Sub(SendBranchSuffix)
}

func (h Home) PrivateReceive() git.Address {
	return h.Private.Sub(ReceiveBranchSuffix)
}

var RootNS = ns.NS{}

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, ProtocolName+": "+msg)
}

const (
	ProtocolName           = "twitter4git" //XXX
	ProtocolVersion        = "0.0.1"
	PostDir                = "post"
	PostFilenameTimeFormat = "20060102-150405"
	SendBranchSuffix       = "main"
	ReceiveBranchSuffix    = "timeline"
	RawExt                 = "raw"
	MetaExt                = "meta.json"
)

type LocalID string // YYYYMMDD-HHMMSS-SHA256CONTENT-NONCE

func (x LocalID) String() string {
	return string(x)
}

func PostNS(by Handle, t time.Time, content string) (ns.NS, LocalID) {
	localID := PostFilebase(by, t, content)
	return RootNS.Join(ns.NS{PostDir, localID.String()}), localID
}

// PostFilebase returns a filename of the form YYYYMMDD-HHMMSS-SHA256CONTENT-NONCE
func PostFilebase(by Handle, t time.Time, content string) LocalID {
	return LocalID(
		t.UTC().Format(PostFilenameTimeFormat) +
			"-" + ContentHash(content) +
			"-" + ContentHash(by.String()) +
			"-" + ContentHash(Nonce()))
}

func ContentHash(content string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(content)); err != nil {
		panic(err)
	}
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func Nonce() string {
	return strconv.Itoa(int(rand.Int63()))
}

type PostMeta struct {
	By Handle `json:"by"`
}
