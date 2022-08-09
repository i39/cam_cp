package filter

import (
	"bytes"
	"cam_cp/app/watcher"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type Filter interface {
	Run(ctx context.Context) error
	In() watcher.ExChan
	Out() watcher.ExChan
}

type content struct {
	fname string
	ftype string
	fdata []byte
}

// sendPostRequest is used to send an HTTP POST request.
func sendPostRequest(url string, file content) ([]byte, error) {
	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)

	part, err := w.CreateFormFile(file.ftype, filepath.Base(file.fname))
	if err != nil {
		return []byte{}, err
	}
	_, err = part.Write(file.fdata)
	if err != nil {
		return []byte{}, err
	}

	err = w.Close()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	var client = new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return cnt, nil
}
