package http_utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// Content is a struct which contains a file's name, its type and its data.
type Content struct {
	Fname string
	Ftype string
	Fdata []byte
}

// SendGetRequest is used to send an HTTP POST request.
func SendGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

// SendPostRequest is used to send an HTTP POST request.
func SendPostRequest(url string, files ...Content) ([]byte, error) {
	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)

	for _, f := range files {
		part, err := w.CreateFormFile(f.Ftype, filepath.Base(f.Fname))
		if err != nil {
			return []byte{}, err
		}
		_, err = part.Write(f.Fdata)
		if err != nil {
			return []byte{}, err
		}
	}

	err := w.Close()
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
