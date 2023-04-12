package watcher

import (
	"cam_cp/app/frame"
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/go-pkgz/lgr"
	mjpegsrv "github.com/mattn/go-mjpeg"
)

// https://github.com/mattn/go-mjpeg

type HttpMJpeg struct {
	Url         string
	decoder     *mjpegsrv.Decoder
	frameBuffer int
	// Interval time.Duration
}

func NewHttpMJpeg(url string, frameBuffer int) (w HttpMJpeg, err error) {
	d, err := mjpegsrv.NewDecoderFromURL(url)
	if err != nil {
		return w, err
	}
	return HttpMJpeg{
		Url:         url,
		decoder:     d,
		frameBuffer: frameBuffer,
	}, nil
}

func (m HttpMJpeg) Watch(ctx context.Context, frames chan<- []frame.Frame) error {
	if frames == nil {
		return fmt.Errorf("frames channel is nil")
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			frameSet := make([]frame.Frame, m.frameBuffer)
			realFrames := 0
			for i := 0; i < m.frameBuffer; i++ {
				img, err := m.decoder.DecodeRaw()
				if err != nil {
					log.Printf("[ERROR] MJPEG watcher error:: %v", err)
					//if error start new decoder
					m.decoder, err = mjpegsrv.NewDecoderFromURL(m.Url)
					if err != nil {
						return err
					}
					continue
				}
				timeStamp := time.Now().UTC().UnixNano()
				nameFromTimeStamp := strconv.FormatInt(timeStamp, 10) + ".jpg"
				f := frame.Frame{Name: nameFromTimeStamp, Data: img, Timestamp: timeStamp}
				frameSet[i] = f
				realFrames++
			}
			if realFrames > 0 {
				frames <- frameSet[:realFrames]
			}
		}
	}
}
