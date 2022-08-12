package sender

//https://github.com/NicoNex/echotron/

import (
	"cam_cp/app/frame"
)

type Sender interface {
	Send(frames []frame.Frame) error
}
