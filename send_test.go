package sendcloud

import (
	"flag"
	"testing"
)

// fake email
type mail struct {
	from    string
	to      []string
	subject string
	html    string
}

func (m *mail) From() string               { return m.from }
func (m *mail) To() []string               { return m.to }
func (m *mail) Cc() []string               { return nil }
func (m *mail) Bcc() []string              { return nil }
func (m *mail) ReplyTo() string            { return "" }
func (m *mail) Subject() string            { return m.subject }
func (m *mail) Html() string               { return m.html }
func (m *mail) Headers() map[string]string { return nil }

var (
	domain = flag.String("domain", "", "Sendcloud domain")
	user   = flag.String("username", "", "Sendcloud username")
	pswd   = flag.String("password", "", "Sendcloud password")
	from   = flag.String("from", "", "From address")
	to     = flag.String("to", "", "To address")
	sc     *Sendcloud
)

func init() {
	flag.Parse()
	sc = New()
	sc.AddDomain(*domain, *user, *pswd)
}

func TestSend(t *testing.T) {
	email := &mail{
		from:    *from,
		to:      []string{*to},
		subject: "SendCloud test mail",
		html:    "SendCloud test mail body",
	}
	mailId, err := sc.Send(email)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mail-id = %s", mailId)
}
