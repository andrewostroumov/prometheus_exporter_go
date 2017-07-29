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
	gpuTemp0 = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gpu_temperature",
			Help: "GPU Temperature",
			ConstLabels: map[string]string{ "index": "0", "name": "RX470", "type": "AMD" },
		},
	)

	gpuTemp1 = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gpu_temperature",
			Help: "GPU Temperature",
			ConstLabels: map[string]string{ "index": "1", "name": "RX480", "type": "AMD" },
		},
	)

	gpuFanSpeed0 = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gpu_fan_speed",
			Help: "GPU Fan Speed",
			ConstLabels: map[string]string{ "index": "0", "name": "RX470", "type": "AMD" },
		},
	)

	gpuFanSpeed1 = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gpu_fan_speed",
			Help: "GPU Fan Speed",
			ConstLabels: map[string]string{ "index": "1", "name": "RX480", "type": "AMD" },
		},
	)

)


func init() {
	prometheus.MustRegister(gpuTemp0)
	prometheus.MustRegister(gpuTemp1)
	prometheus.MustRegister(gpuFanSpeed0)
	prometheus.MustRegister(gpuFanSpeed1)
}

func main() {
	flag.Parse()


	go func() {
		for {
			v := rand.Float64()
			gpuTemp0.Set(v)
			gpuTemp1.Set(v)
			gpuFanSpeed0.Set(v)
			gpuFanSpeed1.Set(v)
			time.Sleep(time.Duration(10000) * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
