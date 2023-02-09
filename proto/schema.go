package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

type Handle string // user's public handle, which is a host/path pair

func (h Handle) String() string {
	return string(h)
}

func (h Handle) URL(ctx context.Context) git.URL {
	u, err := url.Parse(string(h))
	must.NoError(ctx, err)
	if u.Port() == "" {
		return git.URL("https://" + filepath.Join(u.Host, u.Path))
	} else {
		return git.URL("https://" + filepath.Join(u.Host+":"+u.Port(), u.Path))
	}
}

func ParseHandle(urlOrHandle string) (Handle, error) {
	u, err := url.Parse(urlOrHandle)
	if err != nil {
		return "", err
	}
	if u.Port() == "" {
		return Handle(filepath.Join(u.Host, u.Path)), nil
	} else {
		return Handle(filepath.Join(u.Host+":"+u.Port(), u.Path)), nil
	}
}

type Home struct {
	Handle     Handle
	PublicURL  git.URL
	PrivateURL git.URL
}

func (h Home) PublicReadOnly(ctx context.Context) git.Address {
	return git.NewAddress(h.Handle.URL(ctx), PublicBranch)
}

func (h Home) PublicReadWrite() git.Address {
	return git.NewAddress(h.PublicURL, PublicBranch)
}

func (h Home) PrivateReadWrite() git.Address {
	return git.NewAddress(h.PrivateURL, PrivateBranch)
}

var RootNS = ns.NS{}

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, ProtocolName+": "+msg)
}

const (
	ProtocolName           = "skrit4git"
	ProtocolVersion        = "0.0.1"
	PostDir                = "post"
	PostFilenameTimeFormat = "20060102-150405"
	PublicBranch           = "main"
	PrivateBranch          = "timeline"
	RawExt                 = "raw"
	MetaExt                = "meta.json"
)

type LocalID string // YYYYMMDD-HHMMSS-SHA256CONTENT-NONCE

func (x LocalID) String() string {
	return string(x)
}

func PostNS(by Handle, t time.Time, content string) (ns.NS, LocalID) {
	localID := PostFilebase(by, t, content)
	year := fmt.Sprintf("%04d", t.UTC().Year())
	month := fmt.Sprintf("%02d", t.UTC().Month())
	day := fmt.Sprintf("%02d", t.UTC().Day())
	return ns.NS{PostDir, year, month, day, localID.String()}, localID
}

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

// PostFilebase returns a filename of the form YYYYMMDD-HHMMSS-SHA256CONTENT-NONCE
func PostFilebase(by Handle, t time.Time, content string) LocalID {
	return LocalID(
		t.UTC().Format(PostFilenameTimeFormat) +
			"-" + ContentHash(content) +
			"-" + ContentHash(by.String()) +
			"-" + ContentHash(Nonce()))
}

func CacheBranch(url git.URL) string {
	return filepath.Join("cache", ContentHash(string(url)))
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

type Following map[Handle]bool
