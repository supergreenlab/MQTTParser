package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"regexp"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	colorTrimExpr = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	msgExpr       = regexp.MustCompile(`([A-Z]) \(([0-9]+)\) ([A-Z]+): @([A-Z]+) ([^$]+)`)
	kvExpr        = regexp.MustCompile(`(([A-Z0-9_a-z]+) ?= ?(-?[A-Z0-9_a-z.]+))+`)

	server   = flag.String("mqtt_server", "tcp://mqtt:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic    = flag.String("mqtt_topic", "#", "Topic to subscribe to")
	qos      = flag.Int("mqtt_qos", 0, "The QoS to subscribe to messages at")
	clientid = flag.String("mqtt_clientid", "", "A clientid for the connection")
	username = flag.String("mqtt_username", "", "A username to authenticate to the MQTT server")
	password = flag.String("mqtt_password", "", "Password to match username")
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	rl := newRawLog(message.Topic(), string(message.Payload()))
	setLastSeen(rl.Id)

	if msgExpr.Match([]byte(rl.Payload)) {
		l := newLog(rl)

		if kvExpr.Match([]byte(l.Msg)) {
			kvl := newKeyValueLog(l)

			fmt.Println("kvl: ")
			fmt.Println(kvl)
			sendRedisKeyValueLog(kvl)
			sendPromKeyValueLog(kvl)
			indexLog(kvl)
		} else {
			fmt.Println("l: ")
			fmt.Println(l)
			indexLog(l)
		}
	} else {
		fmt.Println("rl: ")
		fmt.Println(rl)
		indexLog(rl)
	}
}

func init_mqtt() {
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
