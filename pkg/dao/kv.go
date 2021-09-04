package dao

import (
	"github.com/go-redis/redis/v8"
	"mlauth/pkg/conf"
)

func getKv() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.KvAddr,
		Password: conf.KvPassword,
		DB:       conf.KvDb,
	})
}
