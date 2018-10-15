package xmail

import (
	"testing"
	"fmt"
)

func TestSend(t *testing.T) {
	f := (*Filter)(nil)
	fmt.Println(f.Useful())


	evn := Envelope{}
	evn.From.Email = "myusername@yahoo.com"
	evn.From.Showname = "testname"
	to := AddrEdit{Email:"myusername@gmail.com", Showname:"myshowname"}
	evn.To = append(evn.To, to)
	evn.Subject = "test email title 2"
	body, err := MakeBody()
	if err != nil {
		t.Error(err)
	}

	c := SendContent{}
	c.BodyString = body
	c.BodyType = BodyTypeHTML

	if err := Send(evn, c, "asdfVCXZ*", nil); err != nil {
		t.Error(err)
	}
}
