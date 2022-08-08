package dispatcher

import (
	"cam_cp/app/watcher"
	"context"
	log "github.com/go-pkgz/lgr"
	"sync"
)

type Dispatcher interface {
	Run(ctx context.Context) error
}

type Impl struct {
	sync.Mutex
	in  []watcher.ExChan
	out []watcher.ExChan
}

func (d *Impl) AddIn(in watcher.ExChan) {
	d.Lock()
	defer d.Unlock()
	d.in = append(d.in, in)
}

func (d *Impl) AddOut(out watcher.ExChan) {
	d.Lock()
	defer d.Unlock()
	d.out = append(d.out, out)
}

func (d *Impl) Run(ctx context.Context) error {
	log.Printf("[INFO] dispatcher is started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for _, in := range d.in {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ex := <-in:
					for _, out := range d.out {
						out <- ex
					}
				}
			}
		}
	}
}
