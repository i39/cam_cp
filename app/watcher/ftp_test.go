package watcher

import (
	"bytes"
	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestWalkFtp(t *testing.T) {
	//use dlptest for testing ftp server
	//https://dlptest.com/ftp-test/ for detailed information

	createTempFiles(t)

	f := &Ftp{
		Dir:           "/cam_cp/",
		CheckInterval: time.Second * 10,
		Ip:            "192.168.89.236",
		User:          "ftpuser",
		Password:      "vad6Udkh",
	}

	files, err := f.walkFtp()
	if err != nil {
		t.Error(err)
	}

	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}

	for k, _ := range files {
		assert.Equal(t, "/cam_cp/test2/test3/test3.txt", k)
		//assert.Equal(t, "test2", files["/cam_cp/test2/test2.txt"])
		//assert.Equal(t, "test1", files["/cam_cp/test1.txt"])
	}
	deleteTestFiles(t)
}

func createTempFiles(t *testing.T) {
	c, err := ftp.Dial("192.168.89.236:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}

	err = c.Login("ftpuser", "vad6Udkh")
	if err != nil {
		t.Error(err)
	}

	err = c.MakeDir("/cam_cp")
	if err != nil {
		t.Error(err)
	}
	err = c.MakeDir("/cam_cp/test2")
	if err != nil {
		t.Error(err)
	}

	err = c.MakeDir("/cam_cp/test2/test3")
	if err != nil {
		t.Error(err)
	}

	var data *bytes.Buffer
	data = bytes.NewBufferString("test1")
	err = c.Stor("/cam_cp/test1.txt", data)
	if err != nil {
		t.Error(err)
	}

	data = bytes.NewBufferString("test2")
	err = c.Stor("/cam_cp/test2/test2.txt", data)
	if err != nil {
		t.Error(err)
	}

	data = bytes.NewBufferString("test3")
	err = c.Stor("/cam_cp/test2/test3/test3.txt", data)
	if err != nil {
		t.Error(err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}

func deleteTestFiles(t *testing.T) {
	c, err := ftp.Dial("ftp.dlptest.com:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}

	err = c.Login("dlpuser", "rNrKYTX9g7z3RgJRmxWuGHbeu")
	if err != nil {
		t.Error(err)
	}

	err = c.RemoveDirRecur("/cam_cp")
	if err != nil {
		t.Error(err)
	}
	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}
