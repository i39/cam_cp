package watcher

import (
	"cam_cp/app/frame"
	"cam_cp/app/http_utils"
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/go-pkgz/lgr"
)

type HttpJpeg struct {
	Url            string
	CheckInterval  time.Duration
	frameBuffer    []frame.Frame
	framesInBuffer int
}

func NewHttpJpeg(url string, frameBufferLen int, checkInterval int64) (HttpJpeg, error) {
	return HttpJpeg{
		Url:            url,
		CheckInterval:  time.Duration(checkInterval),
		frameBuffer:    make([]frame.Frame, frameBufferLen),
		framesInBuffer: 0,
	}, nil
}

func (h HttpJpeg) Watch(ctx context.Context, frames chan<- []frame.Frame) error {
	if frames == nil {
		return fmt.Errorf("frames channel is nil")
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(h.CheckInterval * time.Millisecond):
			if h.framesInBuffer == len(h.frameBuffer) {
				h.framesInBuffer = 0
				frames <- h.frameBuffer
			}
			img, err := http_utils.SendGetRequest(h.Url)
			if err != nil {
				log.Printf("[ERROR] httpJpeg watcher: %s", err)
				continue
			}
			timeStamp := time.Now().UTC().UnixNano()
			nameFromTimeStamp := strconv.FormatInt(timeStamp, 10) + ".jpg"
			f := frame.Frame{Name: nameFromTimeStamp, Data: img, Timestamp: timeStamp}
			h.frameBuffer[h.framesInBuffer] = f
			h.framesInBuffer++
		}
	}
}
