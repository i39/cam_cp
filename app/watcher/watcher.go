//Package watcher provides a common interfaces for watching
//ftp server, file directory, http url, etc.

package watcher

import (
	"cam_cp/app/frame"
	"context"
)

type Watcher interface {
	Watch(ctx context.Context, frames chan<- []frame.Frame) error
}
