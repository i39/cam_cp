package filter

import (
	"cam_cp/app/frame"
	"cam_cp/app/http_utils"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/go-pkgz/lgr"
)

type Deepstack struct {
	Url        string
	ApiKey     string
	Labels     []string
	Confidence float64
}

type Predictions struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
	YMin       int     `json:"y_min"`
	XMin       int     `json:"x_min"`
	YMax       int     `json:"y_max"`
	XMax       int     `json:"x_max"`
}
type OdResponse struct {
	Ok          bool          `json:"success"`
	Error       string        `json:"error"`
	Duration    int           `json:"duration"`
	Predictions []Predictions `json:"predictions"`
}

func NewDeepstack(url, apiKey string, labels string, confidence float64) (d Deepstack, err error) {
	lbs := strings.Split(labels, ",")
	d = Deepstack{Url: url, ApiKey: apiKey,
		Labels: lbs, Confidence: confidence,
	}
	return d, nil
}

func (f Deepstack) Filter(inFrames []frame.Frame) (outFrames []frame.Frame) {

nextFrame:
	for _, fr := range inFrames {
		predictions, err := f.detect(fr)
		if err != nil {
			log.Printf("[ERROR] deepstack filter: %s", err)
			continue
		}
		for _, p := range predictions {
			for _, l := range f.Labels {
				if p.Label == l && p.Confidence >= f.Confidence {
					outFrames = append(outFrames, fr)
					//go to next frame if found any labeled object
					continue nextFrame
				}
			}
		}

	}
	return outFrames
}

func (f Deepstack) detect(inFrame frame.Frame) (pr []Predictions, err error) {
	var cnt = http_utils.Content{Fname: inFrame.Name, Ftype: "image", Fdata: inFrame.Data}
	var dsRes OdResponse
	res, err := http_utils.SendPostRequest(f.Url, cnt)
	if err != nil {
		return pr, err
	}

	if err = json.Unmarshal(res, &dsRes); err != nil {
		return pr, err
	}

	if !dsRes.Ok {
		return pr, fmt.Errorf("deepstack error: %s", dsRes.Error)
	}
	pr = dsRes.Predictions
	return pr, nil
}

func (f Deepstack) Close() {
	//nothing to close
}
