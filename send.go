package sendcloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Mail interface {
	From() string // mail-from address
	To() []string
	Cc() []string
	Bcc() []string
	ReplyTo() string // reply-to address
	Subject() string
	Html() string               // HTML mail body
	Headers() map[string]string // extra mail headers
}

var EMAIL_DOMAIN_RE = regexp.MustCompile(`[^<>]+<?.+@([^<>]+)>?`)

func (sc *Sendcloud) Send(m Mail) (id string, err error) {
	// extract the sending domain
	match := EMAIL_DOMAIN_RE.FindStringSubmatch(m.From())
	if len(match) != 2 {
		err = fmt.Errorf("invalid From address: %s", m.From())
		return
	}
	domain := match[1]

	d := url.Values{}
	d.Add("resp_email_id", "true")
	d.Add("from", m.From())
	if to := m.To(); len(to) > 0 {
		d.Add("to", strings.Join(to, ";"))
	}
	if cc := m.Cc(); len(cc) > 0 {
		d.Add("cc", strings.Join(cc, ";"))
	}
	if bcc := m.Bcc(); len(bcc) > 0 {
		d.Add("bcc", strings.Join(bcc, ";"))
	}
	if replyto := m.ReplyTo(); replyto != "" {
		d.Add("replyto", replyto)
	}
	d.Add("subject", m.Subject())
	d.Add("html", m.Html())

	headers := m.Headers()
	if headers != nil {
		hb, err := json.Marshal(headers)
		if err != nil {
			return "", err
		}
		d.Add("headers", string(hb))
	}

	body, err := sc.do("mail.send", domain, d)
	if err != nil {
		return
	}

	var reply struct {
		Msg  string   `json:"message"`
		Errs []string `json:"errors"`
		Ids  []string `json:"email_id_list"`
	}
	json.Unmarshal(body, &reply)
	if reply.Msg != "success" {
		if len(reply.Errs) > 0 {
			err = fmt.Errorf("SendCloud error: %s", reply.Errs[0])
		} else {
			err = fmt.Errorf("SendCloud error: unknown")
		}
		return
	}
	if len(reply.Ids) > 0 {
		id = reply.Ids[0]
	}
	return
}
