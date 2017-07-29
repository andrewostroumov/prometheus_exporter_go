package main

import (
	"flag"
	"log"
	"net/http"
	"time"
	//"fmt"
	"io/ioutil"
	"strconv"
	"strings"
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
			stats := BuildAMDGPUStats()
			gpuTemp0.Set(stats[0].Temp)
			gpuTemp1.Set(stats[1].Temp)
			gpuFanSpeed0.Set(stats[0].FanSpeed)
			gpuFanSpeed1.Set(stats[1].FanSpeed)
			time.Sleep(time.Duration(10000) * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type AMDGPUStats struct {
	Temp float64
	FanSpeed float64
}

func BuildAMDGPUStats() []AMDGPUStats {
	tempData0, _ := ioutil.ReadFile("/sys/class/drm/card0/device/hwmon/hwmon0/temp1_input")
	tempString0 := strings.TrimSpace(string(tempData0))
	temp0, _ := strconv.ParseFloat(tempString0, 64)

	pwmData0, _ := ioutil.ReadFile("/sys/class/drm/card0/device/hwmon/hwmon0/pwm1")
	pwmMaxData0, _ := ioutil.ReadFile("/sys/class/drm/card0/device/hwmon/hwmon0/pwm1_max")
	pwmString0 := strings.TrimSpace(string(pwmData0))
	pwmMaxString0 := strings.TrimSpace(string(pwmMaxData0))
	pwm0, _ := strconv.ParseFloat(pwmString0, 64)
	pwmMax0, _ := strconv.ParseFloat(pwmMaxString0, 64)
	pwmPercent0 := pwm0 * 100 / pwmMax0

	tempData1, _ := ioutil.ReadFile("/sys/class/drm/card1/device/hwmon/hwmon1/temp1_input")
	tempString1 := strings.TrimSpace(string(tempData1))
	temp1, _ := strconv.ParseFloat(tempString1, 64)

	pwmData1, _ := ioutil.ReadFile("/sys/class/drm/card1/device/hwmon/hwmon1/pwm1")
	pwmMaxData1, _ := ioutil.ReadFile("/sys/class/drm/card1/device/hwmon/hwmon1/pwm1_max")
	pwmString1 := strings.TrimSpace(string(pwmData1))
	pwmMaxString1 := strings.TrimSpace(string(pwmMaxData1))
	pwm1, _ := strconv.ParseFloat(pwmString1, 64)
	pwmMax1, _ := strconv.ParseFloat(pwmMaxString1, 64)
	pwmPercent1 := pwm1 * 100 / pwmMax1


	return []AMDGPUStats{
		{
			Temp: temp0 / 1000,
			FanSpeed: pwmPercent0,
		},
		{
			Temp: temp1 / 1000,
			FanSpeed: pwmPercent1,
		},
	}
}