package watcher

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"github.com/jlaffaye/ftp"
	"time"
)

type Ftp struct {
	Dir           string
	CheckInterval time.Duration
	Ip            string
	User          string
	Password      string
}

func (f *Ftp) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.CheckInterval):
			log.Printf("[DEBUG] file watcher: %s", f.Dir)
			f, err := f.walkFtp()
			if err != nil {
				log.Printf("[ERROR] file watcher: %s", err)
			}
			log.Printf("[DEBUG] file watcher: %s", f)
		}
	}
}

//walt thought ftp directory

func (f *Ftp) walkFtp() ([]string, error) {
	var e *ftp.Entry
	var files []string
	if i := last(f.Ip, ':'); i < 0 {
		f.Ip += ":21"
	}

	c, err := ftp.Dial(f.Ip, ftp.DialWithTimeout(time.Second*10))
	if err != nil {
		return nil, err
	}
	err = c.Login(f.User, f.Password)
	if err != nil {
		return nil, err
	}

	w := c.Walk(f.Dir)
	for {
		if !w.Next() {
			break
		}
		e = w.Stat()
		if e.Type == ftp.EntryTypeFile {
			files = append(files, w.Path())
		}
	}

	return files, nil
}

func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}
