package sender

import (
	"cam_cp/app/watcher"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/go-pkgz/lgr"
)

type inputMediaPhoto struct {
	Type  string `json:"type"`
	Media string `json:"media"`
}

type Telegram struct {
	Token  string
	ChatId int64
}

func (t *Telegram) Run(ctx context.Context, in watcher.ExChan) (err error) {
	log.Printf("[INFO] telegram sender is started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ex := <-in:
			t.send(ex)
		}
	}
}
func (t *Telegram) send(ex []watcher.Exchange) {
	i10 := len(ex) / 10
	r10 := len(ex) % 10
	for i := 0; i < i10; i++ {
		ex10 := ex[i*10 : (i+1)*10]
		err := t.sendMediaGroup(ex10)
		if err != nil {
			log.Printf("[ERROR] telegram sender: %s", err)
			return
		}
	}
	err := t.sendMediaGroup(ex[i10*10 : i10*10+r10])
	if err != nil {
		log.Printf("[ERROR] telegram sender: %s", err)
		return
	}
}

// sendMediaGroup send a group of photos or videos as an album.
func (t *Telegram) sendMediaGroup(ex []watcher.Exchange) (err error) {
	var (
		cnt []content
		im  []inputMediaPhoto
		jsn []byte
		res []byte
	)
	var url = fmt.Sprintf(
		"%ssendMediaGroup?chat_id=%d",
		fmt.Sprintf("https://api.telegram.org/bot%s/", t.Token),
		t.ChatId,
	)

	for _, e := range ex {
		cnt = append(cnt, content{e.Name, e.Name, e.Data})
		im = append(im, inputMediaPhoto{Type: "photo", Media: e.Name})
	}
	jsn, err = json.Marshal(im)
	url = fmt.Sprintf("%s&media=%s", url, jsn)
	if len(cnt) > 0 {
		res, err = sendPostRequest(url, cnt...)
		if err != nil {
			log.Printf("[ERROR] telegram sender: %s, response is: %s", err, res)
			return err
		}
	}

	return err
}
