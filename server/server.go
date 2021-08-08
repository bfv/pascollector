package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/bfv/pascollector/types"
)

var wg sync.WaitGroup
var Config types.ConfigFile

func Start(cfg types.ConfigFile) {

	Config = cfg

	wg.Add(1)

	go startListener()
	go dataCollector()
	go dataSender()

	wg.Wait()

	fmt.Println("server stopped")
}

func Stop() {
	defer wg.Done()
	fmt.Println("stopping server...")
	fmt.Println("server stopped")
}

func dataCollector() {

	ticker := time.NewTicker(time.Duration(Config.CollectInterval) * time.Second)

	for {
		<-ticker.C
		CollectData()
	}
}

func dataSender() {

	ticker := time.NewTicker(time.Duration(Config.SendInterval) * time.Second)

	for {
		<-ticker.C
		SendData(Config)
	}
}
