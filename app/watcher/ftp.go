package watcher

import (
	"cam_cp/app/frame"
	"context"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"net"
	"time"
)

type Ftp struct {
	Dir           string
	CheckInterval time.Duration
	Ip            string
	User          string
	Password      string
}

func NewFtp(ip string, dir string, user string,
	password string, checkInterval time.Duration) (f *Ftp, err error) {

	if r := net.ParseIP(ip); r == nil {
		return nil, fmt.Errorf("invalid ip: %s", ip)
	}

	if i := last(ip, ':'); i < 0 {
		ip += ":21"
	}

	f = &Ftp{
		Dir:           dir,
		Ip:            ip,
		User:          user,
		Password:      password,
		CheckInterval: checkInterval,
	}
	return f, nil
}

func (f *Ftp) Watch(ctx context.Context, frames chan<- []frame.Frame) error {
	if frames == nil {
		return fmt.Errorf("frames channel is nil")
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.CheckInterval):
			files, err := f.walkFtp()
			if err != nil {
				log.Printf("[ERROR] ftp watcher: %s", err)
				continue
			}
			frames <- files

		}
	}
}

//walt thought ftp directory

func (f *Ftp) walkFtp() (files []frame.Frame, err error) {
	var e *ftp.Entry
	var r *ftp.Response

	c, err := ftp.Dial(f.Ip, ftp.DialWithTimeout(time.Second*10))
	if err != nil {
		return files, err
	}
	err = c.Login(f.User, f.Password)
	if err != nil {
		return files, err
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
				return files, err
			}
			var b []byte
			b, err = ioutil.ReadAll(r)
			if err != nil {
				return files, err
			}
			files = append(files, frame.Frame{Name: w.Path(), Data: b})
			err = r.Close()
			if err != nil {
				return files, err
			}
			err = c.Delete(w.Path())
			if err != nil {
				log.Printf("[ERROR] ftp watcher delete error: %s", err)
			}
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
