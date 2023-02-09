package proto

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

func Sync(
	ctx context.Context,
	home Home,
) git.Change[bool] {

	cloned := git.CloneAll(ctx, home.PrivateReceive())
	chg := SyncLocal(ctx, home, cloned)
	cloned.Push(ctx)
	return chg
}

func SyncLocal(
	ctx context.Context,
	home Home,
	clone git.Cloned,
) git.Change[bool] {

	// XXX: read following

	git.EmbedOnBranch(
		ctx,
		clone.Repo(),
		XXX,
	)

	return XXX
}
