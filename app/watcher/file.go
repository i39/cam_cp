package watcher

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"time"
)

type File struct {
	Dir           string
	CheckInterval time.Duration
}

func (s *File) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.CheckInterval):
			log.Printf("[DEBUG] file watcher: %s", s.Dir)
		}
	}
}
