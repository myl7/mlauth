package dao

import (
	"context"
	"fmt"
	"mlauth/pkg/conf"
	"strconv"
	"time"
)

func SetUserActiveEmail(uid int, code string) error {
	kv := getKv()
	k := "user-active-email/" + code
	err := kv.Set(context.Background(), k, uid, time.Duration(conf.UserActiveEmailAge)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserActiveEmail(code string) (int, error) {
	kv := getKv()
	k := "user-active-email/" + code
	v, err := kv.GetDel(context.Background(), k).Result()
	if err != nil {
		return 0, err
	}

	d, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return d, nil
}

func SetEmailRetry(sub string, uid int) error {
	kv := getKv()
	k := fmt.Sprintf("email-retry/%s/%d", sub, uid)
	err := kv.Set(context.Background(), k, "1", time.Duration(conf.UserActiveEmailRetryInterval)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func CheckEmailRetry(sub string, uid int) bool {
	kv := getKv()
	k := fmt.Sprintf("email-retry/%s/%d", sub, uid)
	err := kv.Get(context.Background(), k).Err()
	return err != nil
}
