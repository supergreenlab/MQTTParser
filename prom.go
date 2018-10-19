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

func sendPromKeyValueLog(kvl KeyValueLog) {
	for k, v := range kvl.Kvi {
		r, ok := registered[k]
		if ok == false {
			r = &PromRegistered{
				summary: prometheus.NewSummaryVec(
					prometheus.SummaryOpts{Name: fmt.Sprintf("s_%s", k)},
					[]string{"id", "module"},
				),
				gauge: prometheus.NewGaugeVec(
					prometheus.GaugeOpts{Name: fmt.Sprintf("g_%s", k)},
					[]string{"id", "module"},
				),
			}
			prometheus.MustRegister(r.summary)
			prometheus.MustRegister(r.gauge)
			registered[k] = r
		}
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
