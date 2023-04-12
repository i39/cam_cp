package filter

import (
	"cam_cp/app/frame"
	"strings"

	"github.com/LdDl/go-darknet"
	log "github.com/go-pkgz/lgr"
)

type Yolo struct {
	Config      string
	Weights     string
	Labels      []string
	Probability float32
	n           darknet.YOLONetwork
	objects     *objectsList
}

func NewYolo(gpuIndex int, config, weights, labels string, threshold, probability float32) (*Yolo, error) {
	var y Yolo
	lbs := strings.Split(labels, ",")

	n := darknet.YOLONetwork{
		GPUDeviceIndex:           gpuIndex,
		NetworkConfigurationFile: config,
		WeightsFile:              weights,
		Threshold:                threshold,
	}

	if err := n.Init(); err != nil {
		return nil, err
	}

	y = Yolo{
		Config:      config,
		Weights:     weights,
		Labels:      lbs,
		Probability: probability,
		n:           n,
		objects:     newObjectsList(),
	}
	return &y, nil
}

func (y *Yolo) Close() {
	err := y.n.Close()
	if err != nil {
		log.Printf("[ERROR] can't close yolo network, %v", err)
	}
}

func (y *Yolo) Filter(inFrames []frame.Frame) (outFrames []frame.Frame) {

nextFrame:
	for _, fr := range inFrames {
		dr, err := y.detect(fr)
		if err != nil {
			log.Printf("[ERROR] can't detect frame, %v", err)
			continue
		}
		noDetectedObjects := true
		for _, d := range dr.Detections {
			for i := range d.ClassIDs {
				for _, l := range y.Labels {
					if d.ClassNames[i] == l && d.Probabilities[i] >= y.Probability {
						object := object{
							ClassID:              d.ClassIDs[i],
							ClassProbability:     d.Probabilities[i],
							BoundingBox:          d.BoundingBox,
							DetectionTime:        fr.Timestamp,
							DetectionTimeChanged: false,
						}
						noDetectedObjects = false
						//if no objects in slice(first frame), append it
						if y.objects.len() == 0 {
							y.objects.add(object)
							outFrames = append(outFrames, fr)
							//one object of any type enough for one frame
							continue nextFrame
						}
						//if no object in slice, append it and add to outFrames
						if !y.objects.find(object) {
							y.objects.add(object)
							outFrames = append(outFrames, fr)
							continue nextFrame
						}
						//check objects with DetectionTimeChanged flag equal false
						if y.objects.detectionTimeNotChanged() {
							outFrames = append(outFrames, fr)

						}

					}
				}
			}
		}
		//clear objectsList if no object detected on frame
		if noDetectedObjects {
			y.objects.reset()
		}
	}

	return outFrames
}

func (y *Yolo) detect(fr frame.Frame) (dt *darknet.DetectionResult, err error) {
	//encode bytes to image
	src, err := fr.Image()
	if err != nil {
		return dt, err
	}
	imgDarknet, err := darknet.Image2Float32(src)
	if err != nil {
		return dt, err
	}
	defer imgDarknet.Close()

	dt, err = y.n.Detect(imgDarknet)
	if err != nil {
		return dt, err
	}

	return dt, nil
}
