package xmail

import (
	"github.com/badoux/checkmail"
	"github.com/smcduck/xdsa/xstring"
)

type Address struct {
	email string
	loginname string
	showname  string
}

func Validate(email string) error {
	return checkmail.ValidateFormat(email)
}

func GetLoginName(email string) (string, error) {
	addr, err := NewAddress(email, "")
	if err != nil {
		return "", err
	}
	return addr.LoginName(), nil
}

func NewAddress(email, showname string) (*Address, error) {
	if err := Validate(email); err != nil {
		return nil, err
	}
	atidx := xstring.IndexAfter(email, "@", 0)
	loginname := email[0:atidx]
	return &Address{email:email, showname:showname, loginname:loginname}, nil
}

func (a *Address) Email() string {
	return a.email
}

func (a *Address) EmailReplaceLoginNameTail(replace string) string {
	return xstring.TrySubstrLenAscii(a.loginname, 0, len(a.loginname) - len(replace)) + replace + "@" + a.Host()
}

func (a *Address) LoginName() string {
	return a.loginname
}

func (a *Address) ShowName() string {
	return a.showname
}

func (a *Address) Host() string {
	return xstring.LastSubstrByLenAscii(a.email, len(a.email) - (len(a.loginname) + len("@")))
}