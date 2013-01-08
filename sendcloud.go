package sendcloud

import (
	"fmt"
	"io/ioutil"
	"log"
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

type response struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
	Ids     []string `json:"email_id_list"`
}

func (sc Sendcloud) send(email *Email) (id string, err error) {
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

	reply, err := sc.do("mail.send", d)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", reply)
	return
}
