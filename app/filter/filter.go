package filter

import (
	"cam_cp/app/frame"
)

type Filter interface {
	Filter(inFrames []frame.Frame) (outFrames []frame.Frame)
	Close()
}
