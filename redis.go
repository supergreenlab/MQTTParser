package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var (
	r        *redis.Client
	kvserver = flag.String("redis_server", "redis:6379", "Url to the redis instance")
)

func addId(id string) {
	err := r.HSet("pcbs", id, id).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func setLastSeen(id string) {
	key := fmt.Sprintf("%s.last_seen", id)
	err := r.Set(key, time.Now().Unix(), 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func sendRedisKeyValueLog(kvl KeyValueLog) {
	for k, v := range kvl.Kvs {
		key := fmt.Sprintf("%s.%s.%s", kvl.Id, kvl.Module, k)
		err := r.Set(key, v, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	for k, v := range kvl.Kvi {
		key := fmt.Sprintf("%s.%s.%s", kvl.Id, kvl.Module, k)
		err := r.Set(key, v, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func init_redis() {
	r = redis.NewClient(&redis.Options{
		Addr:     *kvserver,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
