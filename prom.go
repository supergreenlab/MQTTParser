package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromRegistered struct {
	summary *prometheus.SummaryVec
	gauge   *prometheus.GaugeVec
}

var registered map[string]*PromRegistered = map[string]*PromRegistered{}

func getPromRegistered(name string) *PromRegistered {
	r, ok := registered[name]
	if ok == false {
		r = &PromRegistered{
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

func sendPromFirstConnect(l Log) {
	r := getPromRegistered("reboots")
	r.gauge.WithLabelValues(l.Id, l.Module).Add(1)
}

func sendPromKeyValueLog(kvl KeyValueLog) {
	for k, v := range kvl.Kvi {
		r := getPromRegistered(k)
		r.summary.WithLabelValues(kvl.Id, kvl.Module).Observe(float64(v))
		r.gauge.WithLabelValues(kvl.Id, kvl.Module).Set(float64(v))
	}
}

func startPromHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init_prom() {
	go startPromHTTP()
}
