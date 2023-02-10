package proto

import (
	"context"
	"time"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Post(
	ctx context.Context,
	home Home,
	content string,
) git.Change[PostID] {

	cloned := git.CloneOne(ctx, home.PublicReadWrite())
	chg := PostLocal(ctx, home, cloned, content)
	cloned.Push(ctx)
	return chg
}

func PostLocal(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	content string,
) git.Change[PostID] {

	chg := PostStageOnly(ctx, home, clone, content)
	Commit(ctx, clone.Tree(), chg.Msg)
	return chg
}

func PostStageOnly(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	content string,
) git.Change[PostID] {

	postNS, localID := PostNS(home.Handle, time.Now(), content)
	meta := PostMeta{By: home.Handle}
	git.StringToFileStage(ctx, clone.Tree(), postNS.Ext(RawExt), content)
	form.ToFile(ctx, clone.Tree().Filesystem, postNS.Ext(MetaExt).Path(), meta)
	return git.Change[PostID]{
		Result: localID,
		Msg:    "post",
	}
}
