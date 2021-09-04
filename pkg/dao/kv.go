package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"mlauth/pkg/conf"
	"strconv"
	"time"
)

func getKv() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.KvAddr,
		Password: conf.KvPassword,
		DB:       conf.KvDb,
	})
}

func setUid(sub string, age int, uid int, code string) error {
	kv := getKv()
	k := sub + "/" + code
	err := kv.Set(context.Background(), k, uid, time.Duration(age)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func getUid(sub string, code string) (int, error) {
	kv := getKv()
	k := sub + "/" + code
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
