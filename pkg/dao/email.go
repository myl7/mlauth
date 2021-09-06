package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"mlauth/pkg/conf"
	"time"
)

func SetUserActiveEmail(uid int, code string) error {
	return setUid("user-active-email", conf.UserActiveEmailAge, uid, code)
}

func GetUserActiveEmail(code string) (int, error) {
	return getUid("user-active-email", code)
}

func SetUserRecoverEmail(uid int, code string) error {
	return setUid("user-recover-email", conf.UserRecoverEmailAge, uid, code)
}

func GetUserRecoverEmail(code string) (int, error) {
	return getUid("user-recover-email", code)
}

type emailEditEmailBody struct {
	Uid   int    `json:"uid"`
	Email string `json:"email"`
}

func SetEmailEditEmail(uid int, email string, code string) error {
	body := emailEditEmailBody{
		Uid:   uid,
		Email: email,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	kv := getKv()
	k := "email-edit-email/" + code
	err = kv.Set(context.Background(), k, b, time.Duration(conf.EmailEditEmailAge)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetEmailEditEmail(code string) (int, string, error) {
	kv := getKv()
	k := "email-edit-email/" + code
	v, err := kv.GetDel(context.Background(), k).Result()
	if err != nil {
		return 0, "", err
	}

	body := emailEditEmailBody{}
	err = json.Unmarshal([]byte(v), &body)
	if err != nil {
		return 0, "", err
	}

	return body.Uid, body.Email, nil
}

func SetEmailRetry(sub string, uid int) error {
	if conf.EmailRetryInterval == 0 {
		return nil
	}

	kv := getKv()
	k := fmt.Sprintf("email-retry/%s/%d", sub, uid)
	err := kv.Set(context.Background(), k, "1", time.Duration(conf.EmailRetryInterval)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func CheckEmailRetry(sub string, uid int) bool {
	if conf.EmailRetryInterval == 0 {
		return true
	}

	kv := getKv()
	k := fmt.Sprintf("email-retry/%s/%d", sub, uid)
	err := kv.Get(context.Background(), k).Err()
	return err != nil
}
