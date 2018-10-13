package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromRegistered struct {
	summary *prometheus.SummaryVec
}

var registered map[string]*PromRegistered = map[string]*PromRegistered{}

func sendPromKeyValueLog(kvl KeyValueLog) {
	for k, v := range kvl.Kvi {
		r, ok := registered[k]
		if ok == false {
			r = &PromRegistered{
				summary: prometheus.NewSummaryVec(
					prometheus.SummaryOpts{Name: k},
					[]string{"id", "module"},
				),
			}
			prometheus.MustRegister(r.summary)
			registered[k] = r
		}
		r.summary.WithLabelValues(kvl.Id, kvl.Module).Observe(float64(v))
	}
}

func startPromHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init_prom() {
	go startPromHTTP()
}
