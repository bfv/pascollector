/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long:  ``,
	// Run:   runServerCommand,
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the data collector",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the data collector",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		stopServer()
	},
}

func init() {

	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)

	rootCmd.AddCommand(serverCmd)

}

func runServerCommand(cmd *cobra.Command, args []string) {

	for _, server := range Config.Servers {
		fmt.Println(server.Name + ": " + server.Url)
	}
}

func startServer() {

	displayConfig()

	wg.Add(1)

	go startListener()
	go dataCollector()
	go dataSender()

	wg.Wait()

	fmt.Println("server stopped")
}

func stopServer() {
	fmt.Println("stopping server...")
	http.Get("http://localhost:" + strconv.Itoa(Config.Port) + "/stop")
}

func startListener() {

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("stop request received")
		wg.Done()
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
		collectData()
	}
}

func collectData() {
	fmt.Println("collect: " + time.Now().String())
}

func dataSender() {

	ticker := time.NewTicker(time.Duration(Config.SendInterval) * time.Second)

	for {
		<-ticker.C
		sendData()
	}
}

func sendData() {
	fmt.Println("send: " + time.Now().String())
}

func displayConfig() {
	fmt.Println("PASMON v0.0.1")
	fmt.Println("  start datacollector")
	fmt.Println("  port: " + strconv.Itoa(Config.Port))
	fmt.Println("  collect interval: " + strconv.Itoa(Config.CollectInterval) + "s")
	fmt.Println("  send interval: " + strconv.Itoa(Config.SendInterval) + "s")
}
