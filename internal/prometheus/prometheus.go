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

package prometheus

import (
	"fmt"
	"log"
	"net/http"

	mqttparser "github.com/SuperGreenLab/MQTTParser/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type promRegistered struct {
	summary *prometheus.SummaryVec
	gauge   *prometheus.GaugeVec
}

var registered map[string]*promRegistered = map[string]*promRegistered{}

func getPromRegistered(name string) *promRegistered {
	r, ok := registered[name]
	if ok == false {
		r = &promRegistered{
			summary: prometheus.NewSummaryVec(
				prometheus.SummaryOpts{Name: fmt.Sprintf("s_%s", name)},
				[]string{"id", "module"},
			),
			gauge: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{Name: fmt.Sprintf("g_%s", name)},
				[]string{"id", "module"},
			),
		}
		prometheus.MustRegister(r.summary)
		prometheus.MustRegister(r.gauge)
		registered[name] = r
	}
	return r
}

// SendPromFirstConnect on boot message
func SendPromFirstConnect(l mqttparser.Log) {
	r := getPromRegistered("reboots")
	r.gauge.WithLabelValues(l.ID, l.Module).Add(1)
}

// SendPromKeyValueLog on KV message
func SendPromKeyValueLog(kvl mqttparser.KeyValueLog) {
	for k, v := range kvl.Kvi {
		r := getPromRegistered(k)
		r.summary.WithLabelValues(kvl.ID, kvl.Module).Observe(float64(v))
		r.gauge.WithLabelValues(kvl.ID, kvl.Module).Set(float64(v))
	}
}

func startPromHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// InitPrometheus inits the prometheus agent server
func InitPrometheus() {
	go startPromHTTP()
}
