package proto

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

func GetTimelinePostByID(
	ctx context.Context,
	home Home,
	postID PostID,
) PostWithMeta {

	cloned := git.CloneOne(ctx, home.TimelineReadOnly())
	return GetTimelinePostByIDLocal(ctx, cloned, postID)
}

func GetTimelinePostByIDLocal(
	ctx context.Context,
	clone git.Cloned,
	postID PostID,
) PostWithMeta {

	postNS := postID.NS()
	meta := git.FromFile[PostMeta](ctx, clone.Tree(), postNS.Ext(MetaExt).Path())
	content := git.FileToString(ctx, clone.Tree(), postNS.Ext(RawExt))
	return PostWithMeta{Content: []byte(content), Meta: meta}
}
