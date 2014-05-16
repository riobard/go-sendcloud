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

//邮件模板替换参数
type Substitution struct {
	To  []string            `json:"to"`
	Sub map[string][]string `json:"sub"`
}

var EMAIL_DOMAIN_RE = regexp.MustCompile(`[^<>]+<?.+@([^<>]+)>?`)

//返回一个新的替换参数
func NewSubstitution() *Substitution {
	substitutionVars := Substitution{}
	substitutionVars.Sub = make(map[string][]string)
	return &substitutionVars
}

//在替换参数里面增加一个收件人
func (s *Substitution) AddTo(to string) {
	s.To = append(s.To, to)
}

//在替换参数里面增加一组替换
func (s *Substitution) AddSub(search, replace string) {
	s.Sub[search] = append(s.Sub[search], replace)
}

func (c *Client) Send(m Mail) (id string, err error) {
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

	body, err := c.do("mail.send", domain, d)
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

//使用WebAPI发送模板
//substitutionVars里面的To的长度应该和Sub下每一个value的长度一致
func (c *Client) SendTemplate(templateName string, subject string, from string, fromName string, substitutionVars *Substitution) (err error) {
	// extract the sending domain
	match := EMAIL_DOMAIN_RE.FindStringSubmatch(from)
	if len(match) != 2 {
		err = fmt.Errorf("invalid From address: %s", from)
		return
	}
	domain := match[1]

	d := url.Values{}
	d.Add("template_invoke_name", templateName)
	d.Add("subject", subject)
	d.Add("from", from)
	d.Add("fromname", fromName)

	substitutionVarsBytes, err := json.Marshal(substitutionVars)
	if err != nil {
		return err
	}
	d.Add("substitution_vars", string(substitutionVarsBytes))

	body, err := c.do("mail.send_template", domain, d)
	if err != nil {
		return
	}
	var reply struct {
		Msg  string   `json:"message"`
		Errs []string `json:"errors"`
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
	return
}
