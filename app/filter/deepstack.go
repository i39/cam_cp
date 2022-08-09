package filter

import (
	"cam_cp/app/watcher"
	"context"
	log "github.com/go-pkgz/lgr"
)

type Deepstack struct {
	in     watcher.ExChan
	out    watcher.ExChan
	Url    string
	ApiKey string
}

func NewDeepstack(url, apiKey string) *Deepstack {
	return &Deepstack{Url: url, ApiKey: apiKey, in: make(watcher.ExChan), out: make(watcher.ExChan)}
}

func (f *Deepstack) Run(ctx context.Context) error {
	log.Printf("[INFO] deepstack filter for url:%s is started", f.Url)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ex := <-f.in:
			f.send(ex)
		}
	}
}

func (f *Deepstack) In() watcher.ExChan {
	return f.in
}

func (f *Deepstack) Out() watcher.ExChan {
	return f.out
}

func (f *Deepstack) send(ex []watcher.ExData) {
	for _, e := range ex {
		err := f.detect(e)
		if err != nil {
			log.Printf("[ERROR] deepstack filter: %s", err)
		}
	}
}

func (f *Deepstack) detect(ex watcher.ExData) error {
	return nil
}
