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
		Ip:            "ftp.dlptest.com",
		User:          "dlpuser",
		Password:      "rNrKYTX9g7z3RgJRmxWuGHbeu",
	}

	files, err := f.walkFtp()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "/cam_cp/test2/test3/test3.txt", files[0])
	assert.Equal(t, "/cam_cp/test2/test2.txt", files[1])
	assert.Equal(t, "/cam_cp/test1.txt", files[2])
	deleteTestFiles(t)
}

func createTempFiles(t *testing.T) {
	c, err := ftp.Dial("ftp.dlptest.com:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}

	err = c.Login("dlpuser", "rNrKYTX9g7z3RgJRmxWuGHbeu")
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
