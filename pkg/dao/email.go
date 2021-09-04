package dao

import (
	"context"
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
