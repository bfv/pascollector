package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bfv/pascollector/types"
)

var wg sync.WaitGroup
var Config types.ConfigFile

func Start(cfg types.ConfigFile) {

	Config = cfg
	chStore := make(chan types.Metric, 10)

	wg.Add(1)

	go startListener()
	go dataCollector(chStore)
	go dataStore(chStore)
	go dataSender()

	wg.Wait()

	fmt.Println("server stopped")
}

func Stop(cfg types.ConfigFile) {

	Config = cfg
	fmt.Println("stopping server...")
	res, err := http.Get("http://localhost:" + strconv.Itoa(Config.Port) + "/stop")
	if err == nil {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println(err)
	}
}

func dataCollector(ch chan types.Metric) {

	ticker := time.NewTicker(time.Duration(Config.CollectInterval) * time.Second)

	for {
		<-ticker.C
		metrics := CollectData()
		for _, metric := range metrics {
			ch <- metric
		}
	}
}

func dataSender() {

	ticker := time.NewTicker(time.Duration(Config.SendInterval) * time.Second)

	for {
		<-ticker.C
		SendData(Config)
	}
}
