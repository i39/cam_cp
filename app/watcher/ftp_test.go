package watcher

import (
	"bytes"
	"cam_cp/app/frame"
	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

//use dlptest for testing ftp server
//https://dlptest.com/ftp-test/ for detailed information

const FtpTestUrl = "ftp.dlptest.com:21"
const FtpTestDir = "/cam_cp"
const FtpTestUser = "dlpuser"
const FtpTestPassword = "rNrKYTX9g7z3RgJRmxWuGHbeu"

func TestWalkFtp(t *testing.T) {

	createTempFiles(t)

	f := &Ftp{
		Dir:           FtpTestDir,
		CheckInterval: time.Second * 10,
		Ip:            FtpTestUrl,
		User:          FtpTestUser,
		Password:      FtpTestPassword,
	}

	fl := []frame.Frame{
		{
			Name: FtpTestDir + "/test2/test3/test3.txt",
			Data: []byte("test3"),
		},
		{
			Name: FtpTestDir + "/test2/test2.txt",
			Data: []byte("test2"),
		},

		{
			Name: FtpTestDir + "/test1.txt",
			Data: []byte("test1"),
		},
	}

	files, err := f.walkFtp()
	if err != nil {
		t.Error(err)
	}
	deleteTestFiles(t)

	assert.Equal(t, fl, files, "expected %v, got %v", fl, files)

}

func createTempFiles(t *testing.T) {
	c, err := ftp.Dial(FtpTestUrl, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}

	err = c.Login(FtpTestUser, FtpTestPassword)
	if err != nil {
		t.Error(err)
	}

	err = c.MakeDir(FtpTestDir)
	if err != nil {
		t.Error(err)
	}
	err = c.MakeDir(FtpTestDir + "/test2")
	if err != nil {
		t.Error(err)
	}

	err = c.MakeDir(FtpTestDir + "/test2/test3")
	if err != nil {
		t.Error(err)
	}

	var data *bytes.Buffer
	data = bytes.NewBufferString("test1")
	err = c.Stor(FtpTestDir+"/test1.txt", data)
	if err != nil {
		t.Error(err)
	}

	data = bytes.NewBufferString("test2")
	err = c.Stor(FtpTestDir+"/test2/test2.txt", data)
	if err != nil {
		t.Error(err)
	}

	data = bytes.NewBufferString("test3")
	err = c.Stor(FtpTestDir+"/test2/test3/test3.txt", data)
	if err != nil {
		t.Error(err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}

func deleteTestFiles(t *testing.T) {
	c, err := ftp.Dial(FtpTestUrl, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}

	err = c.Login(FtpTestUser, FtpTestPassword)
	if err != nil {
		t.Error(err)
	}

	err = c.RemoveDirRecur(FtpTestDir)
	if err != nil {
		t.Error(err)
	}
	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}
