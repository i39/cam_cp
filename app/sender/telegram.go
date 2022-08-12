package sender

import (
	"cam_cp/app/frame"
	"cam_cp/app/http_utils"
	"encoding/json"
	"fmt"
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
}

func NewTelegram(token string, chatId int64) *Telegram {
	return &Telegram{
		Token:  token,
		ChatId: chatId,
	}
}

func (t *Telegram) Send(frames []frame.Frame) (err error) {
	i10 := len(frames) / 10
	r10 := len(frames) % 10
	for i := 0; i < i10; i++ {
		ex10 := frames[i*10 : (i+1)*10]
		err := t.sendMediaGroup(ex10)
		if err != nil {
			return err
		}
	}
	err = t.sendMediaGroup(frames[i10*10 : i10*10+r10])
	if err != nil {
		return err
	}
	return nil
}

// sendMediaGroup send a group of photos or videos as an album.
func (t *Telegram) sendMediaGroup(fr []frame.Frame) (err error) {
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

	for _, e := range fr {
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
