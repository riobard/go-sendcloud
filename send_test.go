package sendcloud

import (
	"flag"
	"testing"
)

var user = flag.String("user", "", "Sendcloud user")
var key = flag.String("key", "", "Sendcloud key")
var from = flag.String("from", "", "From address")
var to = flag.String("to", "", "To address")
var sc *Sendcloud

func init() {
	flag.Parse()
	sc = New(*user, *key)
}

func TestSend(t *testing.T) {
	email := &Email{
		From:    *from,
		To:      []string{*to},
		Subject: "SendCloud test mail",
		Html:    "SendCloud test mail body",
	}
	sc.send(email)
}
