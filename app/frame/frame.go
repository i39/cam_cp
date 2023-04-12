package frame

import (
	"bytes"
	"image"
	"image/jpeg"
)

type Frame struct {
	Name      string
	Data      []byte
	Timestamp int64
}

func (f Frame) Image() (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(f.Data))
}
