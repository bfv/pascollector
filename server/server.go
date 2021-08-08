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

	wg.Add(1)

	go startListener()
	go dataCollector()
	go dataSender()

	wg.Wait()

	fmt.Println("server stopped")
}

func Stop() {
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

func startListener() {

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		fmt.Println("stop request received")
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET status")
	})

	http.ListenAndServe(":"+strconv.Itoa(Config.Port), nil)
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
