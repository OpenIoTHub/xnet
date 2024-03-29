package xmail

import (
	"testing"
)

func TestSend(t *testing.T) {
	evn := Envelope{}
	evn.From.Email = "user@yahoo.com"
	evn.From.Showname = "testname"
	to := AddrEdit{Email:"user@qq.com", Showname:"myshowname"}
	evn.To = append(evn.To, to)
	evn.Subject = "test email title"
	body, err := MakeBody()
	if err != nil {
		t.Error(err)
	}

	content := SendContent{}
	content.BodyString = body
	content.BodyType = BodyTypeHTML

	if err := Send(evn, content, "mypwd", nil); err != nil {
		t.Error(err)
	}
}
