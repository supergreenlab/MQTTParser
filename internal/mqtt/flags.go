package mqtt

import (
	"github.com/spf13/pflag"
)

var (
	server   = pflag.String("mqtt_server", "tcp://mqtt:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic    = pflag.String("mqtt_topic", "#", "Topic to subscribe to")
	qos      = pflag.Int("mqtt_qos", 0, "The QoS to subscribe to messages at")
	clientid = pflag.String("mqtt_clientid", "", "A clientid for the connection")
	username = pflag.String("mqtt_username", "", "A username to authenticate to the MQTT server")
	password = pflag.String("mqtt_password", "", "Password to match username")
)
