package sender

import (
	"cam_cp/app/http_utils"
	"cam_cp/app/watcher"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/go-pkgz/lgr"
)

type tgResponse struct {
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Ok          bool   `json:"ok"`
}

type inputMediaPhoto struct {
	Type  string `json:"type"`
	Media string `json:"media"`
}

type Telegram struct {
	Token  string
	ChatId int64
	in     watcher.ExChan
}

func NewTelegram(token string, chatId int64) *Telegram {
	return &Telegram{
		Token:  token,
		ChatId: chatId,
		in:     make(watcher.ExChan),
	}
}

func (t *Telegram) Run(ctx context.Context) (err error) {
	log.Printf("[INFO] telegram sender is started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ex := <-t.in:
			err := t.send(ex)
			if err != nil {
				log.Printf("[ERROR] telegram sender: %s", err)
			}
		}
	}
}
func (t *Telegram) send(ex []watcher.ExData) (err error) {
	i10 := len(ex) / 10
	r10 := len(ex) % 10
	for i := 0; i < i10; i++ {
		ex10 := ex[i*10 : (i+1)*10]
		err := t.sendMediaGroup(ex10)
		if err != nil {
			return err
		}
	}
	err = t.sendMediaGroup(ex[i10*10 : i10*10+r10])
	if err != nil {
		return err
	}
	return nil
}

// sendMediaGroup send a group of photos or videos as an album.
func (t *Telegram) sendMediaGroup(ex []watcher.ExData) (err error) {
	var (
		cnt   []http_utils.Content
		im    []inputMediaPhoto
		jsn   []byte
		res   []byte
		tgres tgResponse
	)
	var url = fmt.Sprintf(
		"%ssendMediaGroup?chat_id=%d",
		fmt.Sprintf("https://api.telegram.org/bot%s/", t.Token),
		t.ChatId,
	)

	for _, e := range ex {
		cnt = append(cnt, http_utils.Content{Fname: e.Name, Ftype: e.Name, Fdata: e.Data})
		im = append(im, inputMediaPhoto{Type: "photo",
			Media: fmt.Sprintf("attach://%s", e.Name)})
	}
	if len(cnt) == 0 {
		return fmt.Errorf("no content to send")
	}
	jsn, err = json.Marshal(im)
	url = fmt.Sprintf("%s&media=%s", url, jsn)

	res, err = http_utils.SendPostRequest(url, cnt...)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(res, &tgres); err != nil {
		return err
	}

	if !tgres.Ok {
		return fmt.Errorf("telegram error: %s", tgres.Description)
	}

	return err
}
