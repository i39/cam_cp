package sender

import (
	"cam_cp/app/frame"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Dir string
}

func NewFile(dir string) (f *File, err error) {
	//check if directory exist, if not create it
	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return f, err
		}
	}
	f = &File{Dir: dir}

	return f, nil
}

// Send sends frames to file system
func (f *File) Send(frames []frame.Frame) (err error) {
	for _, fr := range frames {
		err := f.write(fr)
		if err != nil {
			log.Printf("[ERROR] file sender: %s", err)
		}
	}
	return nil
}

func (f *File) write(file frame.Frame) error {
	fBase := filepath.Base(file.Name)
	fDir := filepath.Dir(f.Dir + file.Name)
	//check if directory exist, if not create it
	if _, err := os.Stat(fDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(fDir, 0755)
		if err != nil {
			return err
		}
	}
	newFile, err := os.Create(filepath.Join(fDir, fBase))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("[ERROR] file sender: %s", err)
		}
	}(newFile)
	_, err = newFile.Write(file.Data)
	if err != nil {
		return err
	}

	return nil
}
