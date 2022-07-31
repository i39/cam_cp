package watcher

import (
	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWalkFtp(t *testing.T) {

	f := &Ftp{
		Dir:           "/test/",
		CheckInterval: time.Second * 10,
		Ip:            "192.168.89.236:21",
		User:          "ftpuser",
		Password:      "vad6Udkh",
	}
	c, err := ftp.Dial(f.Ip, ftp.DialWithTimeout(time.Second*10))
	if err != nil {
		t.Fatal(err)
	}
	err = c.Login(f.User, f.Password)
	if err != nil {
		t.Fatal(err)
	}

	files, err := f.walkFtp()
	assert.Equal(t, "/test/test2/test3/test3.txt", files[0])
	assert.Equal(t, "/test/test2/test2.txt", files[1])
	assert.Equal(t, "/test/test1/test1.txt", files[2])
}
