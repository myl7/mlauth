package srv

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/matcornic/hermes/v2"
	"mlauth/pkg/conf"
	"mlauth/pkg/dao"
	"mlauth/pkg/mdl"
	"net/smtp"
	"net/url"
)

func ReqUserActive(u mdl.User) error {
	code, err := genActiveCode(u.Uid)
	if err != nil {
		return err
	}

	err = sendUserActiveEmail(u, code)
	if err != nil {
		return err
	}

	return nil
}

func RunUserActive(code string) error {
	uid, err := dao.GetUserActiveEmail(code)
	if err != nil {
		return err
	}

	uEdit, err := dao.SelectUser(uid)
	if err != nil {
		return err
	}

	uEdit.IsActive = true
	_, err = dao.UpdateUser(uid, uEdit)
	if err != nil {
		return err
	}

	return nil
}

func genActiveCode(uid int) (string, error) {
	d, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	code := d.String()
	err = dao.SetUserActiveEmail(uid, code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func sendUserActiveEmail(u mdl.User, activeCode string) error {
	p, err := url.Parse(conf.SiteHost)
	if err != nil {
		return err
	}

	p.Path = "/emails/active"
	q := p.Query()
	q.Set("active-code", activeCode)
	p.RawQuery = q.Encode()
	link := p.String()

	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "mlauth",
			Link: conf.SiteHost,
		},
	}
	e := hermes.Email{
		Body: hermes.Body{
			Name:   u.DisplayName,
			Intros: []string{"You have received this email because your email address is used in a mlauth account registration."},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to verify your email address and activate the account:",
					Button: hermes.Button{
						Text: "Confirm your registration",
						Link: link,
					},
				},
			},
			Outros: []string{
				"If you did not request a registration, no further action is required on your part.",
			},
			Signature: "mlauth, Copyright © 2021 myl7, source code is licensed under MIT",
		},
	}
	body, err := h.GenerateHTML(e)
	if err != nil {
		return err
	}

	err = sendEmail([]string{u.Email}, body)
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(to []string, body string) error {
	auth := smtp.PlainAuth("", conf.SmtpUsername, conf.SmtpPassword, conf.SmtpHost)
	addr := fmt.Sprintf("%s:%d", conf.SmtpHost, conf.SmtpPort)
	err := smtp.SendMail(addr, auth, conf.SmtpSender, to, []byte(body))
	if err != nil {
		return err
	}

	return nil
}
