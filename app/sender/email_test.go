package sender

import (
	"cam_cp/app/frame"
	"testing"
)

func TestSendEmail(t *testing.T) {
	data, err := base64toPng()
	if err != nil {
		t.Error(err)
	}

	email, err := NewEmail("192.168.4.1", 1025, "",
		"vhod.levo-cc88e36520@flussonic.loca",
		"vhod.levo-cc88e36520@flussonic.loca",
		"Alarm motion", "")
	if err != nil {
		t.Error(err)
	}

	fl := []frame.Frame{
		{
			Name: "/test2/test3/test3.png",
			Data: data,
		},
		{
			Name: "/test2/test2.png",
			Data: data,
		},

		{
			Name: "/test1.png",
			Data: data,
		},
	}

	for f := range fl {
		err = email.send(fl[f])
		if err != nil {
			t.Error(err)
			break
		}
	}
}
