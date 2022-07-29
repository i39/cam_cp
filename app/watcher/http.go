package watcher

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"time"
)

type Http struct {
	Url           string
	CheckInterval time.Duration
}

func (s *Http) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.CheckInterval):
			log.Printf("[DEBUG] http watcher: %s", s.Url)
		}
	}
}
