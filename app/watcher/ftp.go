package watcher

import (
	"context"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
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

func (f *Ftp) walkFtp() (map[string][]byte, error) {
	var e *ftp.Entry
	var r *ftp.Response
	files := make(map[string][]byte)

	if i := last(f.Ip, ':'); i < 0 {
		f.Ip += ":21"
	}

	c, err := ftp.Dial(f.Ip, ftp.DialWithTimeout(time.Second*10), ftp.DialWithDisabledEPSV(true))
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

			r, err = c.Retr(w.Path())
			if err != nil {
				return nil, fmt.Errorf("error reading file %s: %s", w.Path(), err)
			}
			files[w.Path()], err = ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
		}
	}
	defer r.Close()

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
