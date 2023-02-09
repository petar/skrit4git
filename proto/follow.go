package proto

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

func Follow(
	ctx context.Context,
	home Home,
	handle Handle,
) git.Change[bool] {

	cloned := git.CloneOne(ctx, home.PrivateReceive())
	chg := FollowLocal(ctx, home, cloned, handle)
	cloned.Push(ctx)
	return chg
}

func FollowLocal(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	handle Handle,
) git.Change[bool] {

	chg := FollowStageOnly(ctx, home, clone, handle)
	Commit(ctx, clone.Tree(), chg.Msg)
	return chg
}

func FollowStageOnly(
	ctx context.Context,
	home Home,
	clone git.Cloned,
	handle Handle,
) git.Change[bool] {

	following := GetFollowingLocal(ctx, clone)
	already := following[handle]
	following[handle] = true
	followingNS := XXX
	git.ToFileStage(ctx, git.Worktree(ctx, clone.Repo()), followingNS.Path(), following)
	return git.Change[bool]{
		Result: !already,
		Msg:    "follow",
	}
}
