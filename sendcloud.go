package sendcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_ENDPOINT = "https://sendcloud.sohu.com/webapi/"
)

type Sendcloud struct {
	user string
	pswd string
}

func New(user, pswd string) *Sendcloud {
	return &Sendcloud{user: user, pswd: pswd}
}

func (sc *Sendcloud) do(target string, data url.Values) (body []byte, err error) {
	url := fmt.Sprintf("%s%s.json", API_ENDPOINT, target)
	data.Add("api_user", sc.user)
	data.Add("api_key", sc.pswd)
	rsp, err := http.PostForm(url, data)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	if rsp.StatusCode != 200 {
		err = fmt.Errorf("SendCloud error: %d %s", rsp.StatusCode, body)
	}
	return
}

type Mail interface {
	From() string
	To() []string
	Cc() []string
	Bcc() []string
	ReplyTo() string
	Subject() string
	Html() string
}

func (sc *Sendcloud) Send(m Mail) (id string, err error) {
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

	body, err := sc.do("mail.send", d)
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
