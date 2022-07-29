package watcher

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"time"
)

type Ftp struct {
	Dir           string
	CheckInterval time.Duration
	Ip            string
	User          string
	Password      string
}

func (s *Ftp) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.CheckInterval):
			log.Printf("[DEBUG] file watcher: %s", s.Dir)
		}
	}
}
