package watcher

import (
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
	out           ExChan
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
		out:           make(ExChan),
	}
	return f, nil
}

func (f *Ftp) Run(ctx context.Context) error {
	log.Printf("[INFO] ftp watcher for ip:%s , dir:%s is started", f.Ip, f.Dir)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.CheckInterval):
			files, err := f.walkFtp()
			if err != nil {
				log.Printf("[ERROR] ftp watcher: %s", err)
			}
			if len(files) > 0 {
				f.out <- files
			}
		}
	}
}

func (f *Ftp) Out() ExChan {
	return f.out
}

//walt thought ftp directory

func (f *Ftp) walkFtp() ([]ExData, error) {
	var e *ftp.Entry
	var r *ftp.Response
	var files []ExData

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
			files = append(files, ExData{w.Path(), b})
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
