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
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/SuperGreenLab/MQTTParser/internal/prometheus"
	"github.com/SuperGreenLab/MQTTParser/internal/redis"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	colorTrimExpr = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	msgExpr       = regexp.MustCompile(`([A-Z]) \((-?[0-9]+)\) ([A-Z]+): @([A-Z0-9a-z_]+) ([^$]+)`)
	kvExpr        = regexp.MustCompile(`(([A-Z0-9a-z_]+) ?= ?(-?[A-Z0-9_a-z.]+))+`)
	bootExpr      = regexp.MustCompile(`First connect`)
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	if strings.HasSuffix(message.Topic(), "cmd") {
		return
	}
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

			redis.SendRedisKeyValueLog(kvl)
			prometheus.SendPromKeyValueLog(kvl)
			prometheus.SendPromMessageRecieved(true)
		} else {
			redis.SendRedisEventLog(l)
		}
	} else {
		prometheus.SendPromMessageRecieved(false)
		logrus.Warningf("Unknown message: (%s) %s", message.Topic(), string(message.Payload()))
	}
}

// InitMQTT starts the MQTT connection
func InitMQTT() {
	go func() {
		mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
		mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
		connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
		var (
			username = viper.GetString("MQTTUsername")
			password = viper.GetString("MQTTPassword")
		)
		if username != "" {
			connOpts.SetUsername(username)
			if password != "" {
				connOpts.SetPassword(password)
			}
		}
		tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
		connOpts.SetTLSConfig(tlsConfig)

		connOpts.OnConnect = func(c MQTT.Client) {
			if token := c.Subscribe(viper.GetString("MQTTTopic"), byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}
		}

		client := MQTT.NewClient(connOpts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		} else {
			fmt.Printf("Connected to %s\n", *server)
		}
	}()
}
