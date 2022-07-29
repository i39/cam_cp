//Package watcher provides a common interfaces for watching
//ftp server, file directory, http url, etc.

package watcher

import (
	"context"
)

type Watcher interface {
	Run(ctx context.Context) error
}
