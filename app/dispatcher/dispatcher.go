package dispatcher

import (
	"cam_cp/app/watcher"
	"context"
	log "github.com/go-pkgz/lgr"
	"sync"
)

// Dispatcher copy all incoming events to all outgoing channels
type Dispatcher interface {
	Run(ctx context.Context) error
}

// Impl is implementation of Dispatcher
type Impl struct {
	sync.Mutex
	in  []watcher.ExChan
	out []watcher.ExChan
}

// AddIn adds incoming channel to dispatcher
func (d *Impl) AddIn(in ...watcher.ExChan) {
	d.Lock()
	defer d.Unlock()
	for _, i := range in {
		d.in = append(d.in, i)
	}
}

// AddOut adds outgoing channel to dispatcher
func (d *Impl) AddOut(out ...watcher.ExChan) {
	d.Lock()
	defer d.Unlock()
	for _, o := range out {
		d.out = append(d.out, o)
	}
}

// GetIn returns incoming channels
func (d *Impl) GetIn() []watcher.ExChan {
	d.Lock()
	defer d.Unlock()
	return d.in
}

//GetOut returns outgoing channels
func (d *Impl) GetOut() []watcher.ExChan {
	d.Lock()
	defer d.Unlock()
	return d.out
}

// Run starts dispatcher
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
