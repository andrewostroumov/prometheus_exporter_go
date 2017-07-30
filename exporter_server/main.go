package main

import (
	"flag"
	"log"
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

var gpuStats = [2]AMDGPUStats{{ Name: "RX470" }, { Name: "RX480" }}

var tempMetrics [2]prometheus.Gauge
var fanSpeedMetrics [2]prometheus.Gauge

func init() {
	for index, gpu := range gpuStats {
		tempMetrics[index] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gpu_temperature",
				Help: "GPU Temperature",
				ConstLabels: map[string]string{ "index": strconv.Itoa(index), "name": gpu.Name },
			},
		)

		fanSpeedMetrics[index] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gpu_fan_speed",
				Help: "GPU Fan Speed",
				ConstLabels: map[string]string{ "index": strconv.Itoa(index), "name": gpu.Name },
			},
		)
		prometheus.MustRegister(tempMetrics[index])
		prometheus.MustRegister(fanSpeedMetrics[index])
	}
}

func main() {
	flag.Parse()

	go func() {
		for {
			UpdateAMDGPUStats()
			for index, element := range gpuStats {
				tempMetrics[index].Set(element.Temp)
				fanSpeedMetrics[index].Set(element.FanSpeed)
			}
			time.Sleep(time.Duration(10000) * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type AMDGPUStats struct {
	Name string
	Temp float64
	FanSpeed float64
}

func UpdateAMDGPUStats() {
	for index := range gpuStats {
		FillAMDGPUStats(index)
	}
}

func FillAMDGPUStats(index int) {
	tempData, _ := ioutil.ReadFile(fmt.Sprintf("/sys/class/drm/card%v/device/hwmon/hwmon%v/temp1_input", index, index))
	tempString := strings.TrimSpace(string(tempData))
	temp, _ := strconv.ParseFloat(tempString, 64)
	temp = temp / 1000

	pwmData, _ := ioutil.ReadFile(fmt.Sprintf("/sys/class/drm/card%v/device/hwmon/hwmon%v/pwm1", index, index))
	pwmMaxData, _ := ioutil.ReadFile(fmt.Sprintf("/sys/class/drm/card%v/device/hwmon/hwmon%v/pwm1_max", index, index))
	pwmString := strings.TrimSpace(string(pwmData))
	pwmMaxString := strings.TrimSpace(string(pwmMaxData))
	pwm, _ := strconv.ParseFloat(pwmString, 64)
	pwmMax, _ := strconv.ParseFloat(pwmMaxString, 64)
	fanSpeed := pwm * 100 / pwmMax

	gpuStats[index].Temp = temp
	gpuStats[index].FanSpeed = fanSpeed
}