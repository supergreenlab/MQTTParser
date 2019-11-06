package mqtt

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	server   = pflag.String("mqttserver", "tcp://mqtt:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic    = pflag.String("mqtttopic", "#", "Topic to subscribe to")
	qos      = pflag.Int("mqttqos", 0, "The QoS to subscribe to messages at")
	clientid = pflag.String("mqttclientid", "", "A clientid for the connection")
	username = pflag.String("mqttusername", "", "A username to authenticate to the MQTT server")
	password = pflag.String("mqttpassword", "", "Password to match username")
)

func init() {
	viper.SetDefault("MQTTServer", "tcp://mqtt:1883")
	viper.SetDefault("MQTTTopic", "#")
	viper.SetDefault("MQTTQos", 0)
}
