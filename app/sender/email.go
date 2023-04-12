package sender

import (
	"bytes"
	"cam_cp/app/frame"
	"encoding/base64"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"net/smtp"
	"strconv"
	"strings"
)

type Email struct {
	address  string
	password string
	to       string
	from     string
	subject  string
	body     string
}

func NewEmail(host string, port int, password, to, from, subject, body string) (e Email, err error) {
	//TODO check params
	address := host + ":" + strconv.Itoa(port)
	return Email{
		address:  address,
		password: password,
		to:       to,
		from:     from,
		subject:  subject,
		body:     body,
	}, nil
}

func (e Email) Send(frames []frame.Frame) (err error) {
	for _, fr := range frames {
		err := e.send(fr)
		if err != nil {
			log.Printf("[ERROR] email sender: %s", err)
		}
	}
	return nil
}

func (e Email) send(frame frame.Frame) error {
	//i don't use auth right now
	//auth := smtp.PlainAuth("", e.from, e.password, e.host)

	//err := smtp.SendMail(address, nil, e.from, to, data)
	//golang smtp.SendMail() not work with flussonic mail server
	//when i'm write below code for it - it's work
	//main problem is that flussonic returns short codes < 4 symbols
	//and go r.ReadResponce return error because it use parseCodeLine
	//for checking code
	//func parseCodeLine(line string, expectCode int) (code int, continued bool, message string, err error) {
	//	if len(line) < 4 || line[3] != ' ' && line[3] != '-' {
	//		err = ProtocolError("short response: " + line)
	//		return
	//	}

	to := []string{e.to}
	c, err := smtp.Dial(e.address)
	if err != nil {
		return err
	}
	defer c.Close()
	err = c.Hello("darknet")
	if err != nil {
		fmt.Printf("HELLO error")
		return err
	}
	id, err := c.Text.Cmd("MAIL FROM:<%s>", e.from)
	if err != nil {
		return err
	}
	c.Text.StartResponse(id)
	code, _, err := c.Text.ReadResponse(250)
	c.Text.EndResponse(id)

	id, err = c.Text.Cmd("RCPT TO:<%s>", to[0])
	if err != nil {
		return err
	}
	c.Text.StartResponse(id)
	code, _, err = c.Text.ReadResponse(250)
	c.Text.EndResponse(id)

	id, err = cmd(c, "DATA", 354)

	w := c.Text.DotWriter()

	_, err = w.Write(e.buildMail(frame))
	if err != nil {
		return err
	}
	w.Close()

	line, err := c.Text.ReadLine()
	if err != nil {
		return err
	}
	code, err = strconv.Atoi(line[0:3])
	if err != nil {
		return err
	}
	if code != 250 {
		return fmt.Errorf("expected 250, got %d", code)
	}

	id, err = cmd(c, "QUIT", 221)

	if err != nil {
		return err
	}
	return nil
}

func (e Email) buildMail(frame frame.Frame) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", e.from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join([]string{e.to}, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", e.subject))

	boundary := "my-boundary-637"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n",
		boundary))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", e.body))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", frame.Name))
	buf.WriteString(fmt.Sprintf("Content-ID: <%s>\r\n\r\n", frame.Name))

	data := frame.Data

	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	buf.Write(b)
	buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))

	buf.WriteString("--")

	return buf.Bytes()
}

func cmd(c *smtp.Client, cmd string, code int) (id uint, err error) {
	id, err = c.Text.Cmd(cmd)
	if err != nil {
		return id, err
	}
	c.Text.StartResponse(id)
	line, err := c.Text.ReadLine()
	c.Text.EndResponse(id)
	code, err = strconv.Atoi(line[0:3])
	if err != nil {
		return id, err
	}
	if code != code {
		return id, fmt.Errorf("expected 354, got %d", code)
	}
	return id, nil
}
