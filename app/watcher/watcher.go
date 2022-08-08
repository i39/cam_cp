//Package watcher provides a common interfaces for watching
//ftp server, file directory, http url, etc.

package watcher

import (
	"context"
)

type ExData struct {
	Name string
	Data []byte
}

type ExChan chan []ExData

type Watcher interface {
	Run(ctx context.Context) error
	Out() ExChan
}
