package server

import (
	"fmt"

	"github.com/bfv/pascollector/types"
)

func SendData(config types.ConfigFile) {

	chData := make(chan types.Metric)
	chQuit := make(chan bool)

	go getStoredData(chData, chQuit)

forloop:
	for {
		select {
		case metric := <-chData:
			processMetric(metric)
		case <-chQuit:
			break forloop
		}
	}
}

func processMetric(metric types.Metric) {
	// do whatever is necessary with the metrics here
	// for now we just print them
	_, err := fmt.Println("metric:", metric)
	if err == nil {
		deleteMetric(metric.Id)
	}
}
