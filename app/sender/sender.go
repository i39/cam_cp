package sender

//https://github.com/NicoNex/echotron/

import (
	"cam_cp/app/watcher"
	"context"
)

type Sender interface {
	Run(ctx context.Context) error
	In() watcher.ExChan
}
