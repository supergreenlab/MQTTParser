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
	"crypto/tls"
	"flag"
	"fmt"
	"regexp"
	"strings"

	//"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis"
)

var r *redis.Client
var varExpr = regexp.MustCompile(`^[A-Z0-9_]+=[^$]+$`)
var colorTrimExpr = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	payload := colorTrimExpr.ReplaceAll(message.Payload(), []byte(""))
	id := strings.Split(message.Topic(), ".")[0]
	msg := strings.Join(strings.Split(string(payload), ": ")[1:], " ")
	if varExpr.Match([]byte(msg)) {
		varDesc := strings.Split(msg, "=")
		varName := varDesc[0]
		varValue := varDesc[1]
		numValue, err := strconv.Atoi(varValue)
		key := fmt.Sprintf("%s.%s", id, varName)
		if err == nil {
			fmt.Printf("%s=%d\n", key, numValue)
			err := r.Set(key, numValue, 0).Err()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%s=%s\n", key, varValue)
			err := r.Set(key, varValue, 0).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Printf("[%s]: %s\n", id, payload)
	}
}

func main() {
	//MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	hostname, _ := os.Hostname()

	kvserver := flag.String("redis", "localhost:6379", "Url to the redis instance")
	server := flag.String("server", "tcp://broker.supergreenlab.com:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic := flag.String("topic", "#", "Topic to subscribe to")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	r = redis.NewClient(&redis.Options{
		Addr:     *kvserver,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

	<-c
}
