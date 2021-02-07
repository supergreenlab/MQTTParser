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
