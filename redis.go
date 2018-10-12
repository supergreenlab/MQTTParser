package main

import (
	"flag"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	r        *redis.Client
	kvserver = flag.String("redis_server", "redis:6379", "Url to the redis instance")
)

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
