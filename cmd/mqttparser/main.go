/*
 * Copyright (C) 2018  SuperGreenLab <towelie@supergreenlab.com>
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

package main

import (
	"fmt"

	"github.com/SuperGreenLab/MQTTParser/internal/mqtt"
	"github.com/SuperGreenLab/MQTTParser/internal/prometheus"
	"github.com/SuperGreenLab/MQTTParser/internal/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("mqttparser")
	viper.AddConfigPath("/etc/mqttparser")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.SetEnvPrefix("MQTTPARSER")
	viper.AutomaticEnv()

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	redis.InitRedis()
	prometheus.InitPrometheus()
	mqtt.InitMQTT()

	select {}
}
