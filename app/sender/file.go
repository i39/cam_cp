package sender

import (
	"cam_cp/app/frame"
	"errors"
	log "github.com/go-pkgz/lgr"
	"os"
	"path/filepath"
)

type File struct {
	Dir string
}

func NewFile(dir string) (f File, err error) {
	//check if directory exist, if not create it
	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return f, err
		}
	}

	return File{Dir: dir}, nil
}

// Send sends frames to file system
func (f File) Send(frames []frame.Frame) (err error) {
	for _, fr := range frames {
		err := f.write(fr)
		if err != nil {
			log.Printf("[ERROR] file sender: %s", err)
		}
	}
	return nil
}

func (f File) write(file frame.Frame) error {
	fBase := filepath.Base(file.Name)
	log.Printf("[DEBUG] save to file: %s", filepath.Join(f.Dir, fBase))
	newFile, err := os.Create(filepath.Join(f.Dir, fBase))
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
