package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

type HomeAddress git.Address

var RootNS = ns.NS{}

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, "twitter4git: "+msg)
}

const (
	PostDir                = "post"
	PostFilenameTimeFormat = "20060102-150405"
)

func PostNS(t time.Time, content string) ns.NS {
	return ns.NS{PostDir, PostFilebase(t, content)}
}

// PostFilebase returns a filename of the form YYYYMMDD-HHMMSS-sha256:HASH
func PostFilebase(t time.Time, content string) string {
	return t.UTC().Format(PostFilenameTimeFormat) + ContentHash(content)
}

func ContentHash(content string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(content)); err != nil {
		panic(err)
	}
	return "sha256:" + strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}
