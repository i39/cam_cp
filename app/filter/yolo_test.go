package filter

import (
	"cam_cp/app/frame"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYoloDetect(t *testing.T) {

	data, err := base64toPng()
	if err != nil {
		t.Error(err)
	}

	y, err := NewYolo(0, "../../yolov3.cfg",
		"../../yolov3.weights", "person",
		0.25, 90)
	if err != nil {
		t.Error(err)
	}
	fr := frame.Frame{Name: "test", Data: data}
	pr, err := y.detect(fr)
	if err != nil {
		t.Error(err)
	}
	if len(pr.Detections) == 0 {
		t.Error("no objects detected")
	}
	for _, p := range pr.Detections {
		for i := range p.ClassIDs {
			if p.Probabilities[i] > y.Probability {
				fmt.Println("Detection name: ", p.ClassNames[i],
					"Probability: ", p.Probabilities[i])
				assert.Equal(t, p.ClassNames[i], "person", "expected %v, got %v",
					"person", p.ClassNames[i])
			}
		}
	}
}
