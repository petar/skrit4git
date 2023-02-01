package proto

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

func Post(
	ctx context.Context,
	home HomeAddress,
	content string,
) git.ChangeNoResult {

	cloned := git.CloneOrInit(ctx, git.Address(home))
	chg := PostLocal(ctx, cloned, content)
	cloned.Push(ctx)
	return chg
}

func PostLocal(
	ctx context.Context,
	clone git.Cloned,
	content string,
) git.ChangeNoResult {

	chg := PostStageOnly(ctx, clone, content)
	Commit(ctx, clone.Tree(), "post")
	return chg
}

func PostStageOnly(
	ctx context.Context,
	clone git.Cloned,
	content string,
) git.ChangeNoResult {

	XXX
}
