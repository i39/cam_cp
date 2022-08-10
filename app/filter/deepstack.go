package filter

import (
	"cam_cp/app/http_utils"
	"cam_cp/app/watcher"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"strings"
)

type Deepstack struct {
	in         watcher.ExChan
	out        watcher.ExChan
	Url        string
	ApiKey     string
	Labels     []string
	Confidence float64
}

type Predictions struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
	Y_min      int     `json:"y_min"`
	X_min      int     `json:"x_min"`
	Y_max      int     `json:"y_max"`
	X_max      int     `json:"x_max"`
}
type OdResponce struct {
	Ok          bool          `json:"success"`
	Error       string        `json:"error"`
	Duration    int           `json:"duration"`
	Predictions []Predictions `json:"predictions"`
}

func NewDeepstack(url, apiKey string, labels string, confidence float64) (d *Deepstack, err error) {
	lbs := strings.Split(labels, ",")
	d = &Deepstack{Url: url, ApiKey: apiKey,
		Labels: lbs, Confidence: confidence,
		in: make(watcher.ExChan), out: make(watcher.ExChan)}
	return d, nil
}

func (f *Deepstack) Run(ctx context.Context) error {
	log.Printf("[INFO] deepstack filter for url:%s is started", f.Url)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ex := <-f.in:
			f.send(ex)
		}
	}
}

func (f *Deepstack) In() watcher.ExChan {
	return f.in
}

func (f *Deepstack) Out() watcher.ExChan {
	return f.out
}

func (f *Deepstack) send(ex []watcher.ExData) {
	var out []watcher.ExData
	for _, e := range ex {
		predictions, err := f.detect(e)
		if err != nil {
			log.Printf("[ERROR] deepstack filter: %s", err)
			continue
		}
		for _, p := range predictions {
			for _, l := range f.Labels {
				if p.Label == l && p.Confidence > f.Confidence {
					out = append(out, e)
				}
			}
		}

	}
	f.out <- out
}

func (f *Deepstack) detect(ex watcher.ExData) (pr []Predictions, err error) {
	var cnt = http_utils.Content{Fname: ex.Name, Ftype: "image", Fdata: ex.Data}
	var dsRes OdResponce
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
