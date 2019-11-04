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

package mqtt

import (
	"crypto/tls"
	"fmt"
	"regexp"

	"github.com/SuperGreenLab/MQTTParser/internal/prometheus"
	"github.com/SuperGreenLab/MQTTParser/internal/redis"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	colorTrimExpr = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	msgExpr       = regexp.MustCompile(`([A-Z]) \(([0-9]+)\) ([A-Z]+): @([A-Z0-9a-z_]+) ([^$]+)`)
	kvExpr        = regexp.MustCompile(`(([A-Z0-9a-z_]+) ?= ?(-?[A-Z0-9_a-z.]+))+`)
	bootExpr      = regexp.MustCompile(`First connect`)
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	rl := newRawLog(message.Topic(), string(message.Payload()))
	redis.AddID(rl.ID)
	redis.SetLastSeen(rl.ID)

	if msgExpr.Match([]byte(rl.Payload)) {
		l := newLog(rl)
		if bootExpr.Match([]byte(l.Msg)) {
			prometheus.SendPromFirstConnect(l)
		}

		if kvExpr.Match([]byte(l.Msg)) {
			kvl := newKeyValueLog(l)

			fmt.Println("kvl: ")
			fmt.Println(kvl)
			redis.SendRedisKeyValueLog(kvl)
			prometheus.SendPromKeyValueLog(kvl)
		} else {
			fmt.Println("l: ")
			fmt.Println(l)
		}
	} else {
		fmt.Println("rl: ")
		fmt.Println(rl)
	}
}

// InitMQTT starts the MQTT connection
func InitMQTT() {
	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}
}
