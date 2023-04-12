package filter

import (
	"github.com/LdDl/go-darknet"
)

type object struct {
	// ClassID is the class ID of the detected object.
	ClassID int
	// ClassProbability is the probability of the detected object.
	ClassProbability float32
	// BoundingBox is the bounding box of the detected object.
	BoundingBox darknet.BoundingBox
	// DetectionTime is the time it took to detect the object.
	DetectionTime int64
	// DetectionTime changed flag
	DetectionTimeChanged bool
}

type objectsList struct {
	objects []object
}

func newObjectsList() *objectsList {

	return &objectsList{
		objects: []object{},
	}
}

func (o *objectsList) reset() {
	o.objects = []object{}
}

func (o *objectsList) len() int {
	return len(o.objects)
}

func (o *objectsList) add(obj object) {
	o.objects = append(o.objects, obj)
}

func (o *objectsList) find(obj object) bool {
	for i, v := range o.objects {
		(o.objects)[i].DetectionTimeChanged = false
		if v.ClassID == obj.ClassID && bbEqual(v.BoundingBox, obj.BoundingBox) {
			(o.objects)[i].DetectionTime = obj.DetectionTime
			(o.objects)[i].DetectionTimeChanged = true
			return true
		}
	}
	return false
}

// not changed detection time means disappeared object
func (o *objectsList) detectionTimeNotChanged() bool {
	notChanged := false
	//remove object if detection time not changed
	i := 0 // output index
	for _, v := range o.objects {
		if v.DetectionTimeChanged {
			// copy and increment index
			o.objects[i] = v
			i++
		}
	}
	if i < len(o.objects) {
		o.objects = o.objects[:i]
		notChanged = true
	}

	return notChanged
}

// Compare BoundingBoxes
func bbEqual(a, b darknet.BoundingBox) bool {
	//jitter for bounding box
	const jitter = 10
	startPointXDelta := a.StartPoint.X - b.StartPoint.X
	startPointYDelta := a.StartPoint.Y - b.StartPoint.Y
	endPointXDelta := a.EndPoint.X - b.EndPoint.X
	endPointYDelta := a.EndPoint.Y - b.EndPoint.Y
	if startPointXDelta < jitter && startPointYDelta < jitter && endPointXDelta < jitter && endPointYDelta < jitter {
		return true
	}
	return false
}
