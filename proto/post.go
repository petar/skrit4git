package proto

import (
	"context"
	"time"

	"github.com/gov4git/lib4git/git"
)

func Post(
	ctx context.Context,
	home HomeAddress,
	content string,
) git.Change[LocalID] {

	cloned := git.CloneOrInit(ctx, git.Address(home))
	chg := PostLocal(ctx, cloned, content)
	cloned.Push(ctx)
	return chg
}

func PostLocal(
	ctx context.Context,
	clone git.Cloned,
	content string,
) git.Change[LocalID] {

	chg := PostStageOnly(ctx, clone, content)
	Commit(ctx, clone.Tree(), chg.Msg)
	return chg
}

func PostStageOnly(
	ctx context.Context,
	clone git.Cloned,
	content string,
) git.Change[LocalID] {

	postNS, localID := PostNS(time.Now(), content)
	git.StringToFileStage(ctx, clone.Tree(), postNS, content)
	return git.Change[LocalID]{
		Result: localID,
		Msg:    "post",
	}
}
