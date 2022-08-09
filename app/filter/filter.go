package filter

import (
	"cam_cp/app/watcher"
	"context"
)

type Filter interface {
	Run(ctx context.Context) error
	In() watcher.ExChan
	Out() watcher.ExChan
}
