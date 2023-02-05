package proto

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

func Sync(
	ctx context.Context,
	home Home,
) git.Change[bool] {

	cloned := git.CloneOrInit(ctx, home.PrivateReceive())
	chg := SyncLocal(ctx, home, cloned)
	cloned.Push(ctx)
	return chg
}

func SyncLocal(
	ctx context.Context,
	home Home,
	clone git.Cloned,
) git.Change[bool] {

	chg := SyncStageOnly(ctx, home, clone)
	Commit(ctx, clone.Tree(), chg.Msg)
	return chg
}

func SyncStageOnly(
	ctx context.Context,
	home Home,
	clone git.Cloned,
) git.Change[bool] {

	panic("XXX")
}
