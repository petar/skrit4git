package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

// example handle : https://example.com:8080/git/repo.git
// example link to post : skrit4git_https://example.com:8080/git/repo.git?post=20230123231112_abcd_fghi

type Handle struct {
	Scheme string
	Host   string
	Path   string
}

func (h Handle) String() string {
	return string(h.URL())
}

func (h Handle) URL() git.URL {
	return git.URL(h.Scheme + "://" + h.Host + "/" + h.Path)
}

func (h Handle) MarshalJSON() ([]byte, error) {
	s := h.String()
	return json.Marshal(s)
}

func (h *Handle) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	g, err := ParseHandle(s)
	if err != nil {
		return err
	}
	*h = g
	return nil
}

func MustParseHandle(ctx context.Context, urlOrHandle string) Handle {
	h, err := ParseHandle(urlOrHandle)
	must.NoError(ctx, err)
	return h
}

func ParseHandle(s string) (Handle, error) {
	u, err := url.Parse(s)
	if err != nil {
		return Handle{}, err
	}
	if u.Scheme != "https" {
		return Handle{}, fmt.Errorf("handle must be an https url")
	}
	return Handle{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   strings.Trim(u.Path, "/"),
	}, nil
}

type Link struct {
	Handle
	PostID
}

func NewLink(h Handle, id PostID) Link {
	return Link{Handle: h, PostID: id}
}

func (l Link) URL() *url.URL {
	return &url.URL{
		Scheme: ProtocolName + "_" + l.Handle.Scheme,
		Host:   l.Handle.Host,
		Path:   l.Handle.Path + "?post=" + l.PostID.String(),
	}
}

func (l Link) String() string {
	return l.URL().String()
}

func MustParseLink(ctx context.Context, s string) Link {
	l, err := ParseLink(s)
	must.NoError(ctx, err)
	return l
}

func ParseLink(s string) (Link, error) {
	// parse link as url
	u, err := url.Parse(s)
	if err != nil {
		return Link{}, err
	}
	// parse scheme
	if !strings.HasPrefix(u.Scheme, ProtocolName+"_") {
		return Link{}, fmt.Errorf("link scheme not recognized")
	}
	// parse handle
	h := u.Scheme[len(ProtocolName+"_"):] + "://" + u.Host + "/" + strings.TrimLeft(u.Path, "/")
	handle, err := ParseHandle(h)
	if err != nil {
		return Link{}, err
	}
	// parse id
	p := u.Query().Get("post")
	id, err := ParsePostID(p)
	if err != nil {
		return Link{}, err
	}
	return Link{
		Handle: handle,
		PostID: id,
	}, nil
}

type PostID struct {
	Time        time.Time
	ContentHash string
	Nonce       string
}

func NewPostID(t time.Time, content []byte) PostID {
	return PostID{Time: t, ContentHash: ContentHash(content), Nonce: ContentHash(Nonce())}
}

const IDTimeFormat = "20060102150405"

func ParsePostID(s string) (PostID, error) {
	ps := strings.Split(s, "_")
	if len(ps) != 3 {
		return PostID{}, fmt.Errorf("unexpected number of parts in post id")
	}
	t, err := time.Parse(IDTimeFormat, ps[0])
	if err != nil {
		return PostID{}, err
	}
	return PostID{
		Time:        t,
		ContentHash: ps[1],
		Nonce:       ps[2],
	}, nil
}

func (x PostID) String() string {
	t := x.Time.Format(IDTimeFormat)
	return t + "_" + x.ContentHash + "_" + x.Nonce
}

type Home struct {
	Handle       Handle
	TimelineURL  git.URL
	FollowingURL git.URL
}

func (h Home) Link(postID PostID) Link {
	return NewLink(h.Handle, postID)
}

func (h Home) TimelineReadOnly() git.Address {
	return git.NewAddress(h.Handle.URL(), TimelineBranch)
}

func (h Home) TimelineReadWrite() git.Address {
	return git.NewAddress(h.TimelineURL, TimelineBranch)
}

func (h Home) FollowingReadWrite() git.Address {
	return git.NewAddress(h.FollowingURL, FollowingBranch)
}

func NewPostNS(by Handle, t time.Time, content []byte) (ns.NS, PostID) {
	t = t.UTC()
	id := NewPostID(t, content)
	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())
	return ns.NS{PostDir, year, month, day, id.String()}, id
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
