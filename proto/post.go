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
) git.Change[LocalID] {

	cloned := git.CloneOrInit(ctx, git.Address(home.Private))
	chg := PostLocal(ctx, home, cloned, content)
	cloned.Push(ctx)
	return chg
}

func PostLocal(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	content string,
) git.Change[LocalID] {

	chg := PostStageOnly(ctx, home, clone, content)
	Commit(ctx, clone.Tree(), chg.Msg)
	return chg
}

func PostStageOnly(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	content string,
) git.Change[LocalID] {

	postNS, localID := PostNS(time.Now(), content)
	meta := PostMeta{By: home.Handle}
	git.StringToFileStage(ctx, clone.Tree(), postNS.Ext("raw"), content)
	form.ToFile(ctx, clone.Tree().Filesystem, postNS.Ext("meta").Path(), meta)
	return git.Change[LocalID]{
		Result: localID,
		Msg:    "post",
	}
}
