package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func FilterPosts(path ns.NS, _ object.TreeEntry) bool {
	if len(path) != 5 {
		return false
	}
	if path[0] != PostDir {
		return false
	}
	if _, err := strconv.Atoi(path[1]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(path[2]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(path[3]); err != nil {
		return false
	}
	return true
}

var (
	FollowingNS = ns.NS{"following.json"}
	TimelineNS  = ns.NS{}
)

func CacheBranch(url git.URL) string {
	return strings.Join([]string{"cache", ContentHash([]byte(url))}, "/")
}

func ContentHash(content []byte) string {
	h := sha256.New()
	if _, err := h.Write([]byte(content)); err != nil {
		panic(err)
	}
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func Nonce() []byte {
	return []byte(strconv.Itoa(int(rand.Int63())))
}

type PostWithMeta struct {
	Content []byte
	Meta    PostMeta
}

type PostMeta struct {
	By Handle `json:"by"`
}

type Following map[Handle]bool

const (
	ProtocolName           = "skrit4git"
	ProtocolVersion        = "0.0.1"
	PostDir                = "post"
	PostFilenameTimeFormat = "20060102-150405"
	TimelineBranch         = "timeline"
	FollowingBranch        = "following"
	RawExt                 = "raw"
	MetaExt                = "meta.json"
)

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, ProtocolName+": "+msg)
}
