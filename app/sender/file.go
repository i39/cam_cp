package sender

import (
	"cam_cp/app/watcher"
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Dir string
}

func NewFile(dir string) *File {
	return &File{Dir: dir}
}

func (f *File) Run(ctx context.Context, in watcher.ExChan) error {
	log.Printf("[INFO] file sender for dir:%s is started", f.Dir)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ex := <-in:
			f.send(ex)
		}
	}

}

func (f *File) send(ex []watcher.Exchange) {
	for _, e := range ex {
		err := f.write(e)
		if err != nil {
			log.Printf("[ERROR] file sender: %s", err)
		}
	}

}

func (f *File) write(ex watcher.Exchange) error {
	fBase := filepath.Base(ex.Name)
	fDir := filepath.Dir(f.Dir + ex.Name)
	//check if directory exist, if not create it
	if _, err := os.Stat(fDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(fDir, 0755)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(filepath.Join(fDir, fBase))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("[ERROR] file sender: %s", err)
		}
	}(file)
	_, err = file.Write(ex.Data)
	if err != nil {
		return err
	}

	return nil
}
