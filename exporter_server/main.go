package main

import (
	"flag"
	"log"
	"net/http"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

var (
	requestDuraion = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "request_durations_seconds",
			Help: "Request latency distributions.",
		},
	)
)

func main() {
	flag.Parse()

	prometheus.MustRegister(requestDuraion)

	go func() {
		for {
			v := rand.Float64()
			requestDuraion.Set(v)
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
