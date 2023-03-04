package proto

import (
	"context"
)

func FetchLink(
	ctx context.Context,
	home Home,
	link Link,
) PostWithMeta {

	if link.Handle.URL() == home.Handle.URL() {
		return GetTimelinePostByID(ctx, home, link.PostID)
	}

	their := Home{
		Handle:      link.Handle,
		TimelineURL: link.Handle.URL(),
	}
	return GetTimelinePostByID(ctx, their, link.PostID)
}
