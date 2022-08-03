//Package watcher provides a common interfaces for watching
//ftp server, file directory, http url, etc.

package watcher

type Exchange struct {
	Name string
	Data []byte
}

type ExChan chan []Exchange
