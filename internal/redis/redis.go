/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package redis

import (
	"fmt"
	"time"

	mqttparser "github.com/SuperGreenLab/MQTTParser/pkg"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	r *redis.Client
)

// AddID add a controller to the list
func AddID(id string) {
	err := r.HSet("pcbs", id, id).Err()
	if err != nil {
		fmt.Println(err)
	}
}

// SetLastSeen updates the last_seen key for a controller
func SetLastSeen(id string) {
	key := fmt.Sprintf("%s.last_seen", id)
	err := r.Set(key, time.Now().Unix(), 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

// SendRedisKeyValueLog update a key value for a controller
func SendRedisKeyValueLog(kvl mqttparser.KeyValueLog) {
	for k, v := range kvl.Kvs {
		key := fmt.Sprintf("%s.%s.%s", kvl.ID, kvl.Module, k)
		err := r.Set(key, v, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
		r.Publish(fmt.Sprintf("pub.%s", key), v)
	}

	for k, v := range kvl.Kvi {
		key := fmt.Sprintf("%s.%s.%s", kvl.ID, kvl.Module, k)
		err := r.Set(key, v, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
		r.Publish(fmt.Sprintf("pub.%s", key), v)
	}
}

// InitRedis init the redis connection
func InitRedis() {
	r = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("RedisURL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
