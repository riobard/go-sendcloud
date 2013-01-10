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
	return &Sendcloud{user, pswd}
}

func (sc Sendcloud) do(target string, data url.Values) (body []byte, err error) {
	url := fmt.Sprintf("%s%s.json", API_ENDPOINT, target)
	data.Add("api_user", sc.user)
	data.Add("api_key", sc.pswd)
	data.Add("resp_email_id", "true")
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

type Email struct {
	From     string
	FromName string
	To       []string
	Cc       []string
	Bcc      []string
	ReplyTo  string
	Subject  string
	Html     string
}

func (sc Sendcloud) Send(email *Email) (id string, err error) {
	d := url.Values{}
	d.Add("from", email.From)
	if email.FromName != "" {
		d.Add("fromname", email.FromName)
	}
	if len(email.To) > 0 {
		d.Add("to", strings.Join(email.To, ";"))
	}
	if len(email.Cc) > 0 {
		d.Add("cc", strings.Join(email.Cc, ";"))
	}
	if len(email.Bcc) > 0 {
		d.Add("bcc", strings.Join(email.Bcc, ";"))
	}
	if email.ReplyTo != "" {
		d.Add("replyto", email.ReplyTo)
	}
	d.Add("subject", email.Subject)
	d.Add("html", email.Html)

	var reply struct {
		Msg  string   `json:"message"`
		Errs []string `json:"errors"`
		Ids  []string `json:"email_id_list"`
	}

	body, err := sc.do("mail.send", d)
	if err != nil {
		return
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
