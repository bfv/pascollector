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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/bfv/pascollector/types"
	"github.com/shirou/gopsutil/process"
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

var serverListCms = &cobra.Command{
	Use:   "list",
	Short: "List servers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listServers()
	},
}

func init() {

	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)

	rootCmd.AddCommand(serverCmd)

}

func listServers() {
	for _, instance := range Config.PasInstances {
		fmt.Println(instance.Name + ": " + instance.Url)
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
		collectData()
	}
}

func collectData() {

	var ablApps []types.ABLApplication
	//fmt.Println("collect: " + time.Now().String())

	for _, instance := range Config.PasInstances {
		ablApps = getAblApps(instance)

		for idx, ablApp := range ablApps {
			agents := getAgents(instance, ablApp)
			ablApps[idx].Agents = agents

			for idx, agent := range agents {
				agentStats := getAgentStats(instance, ablApp, agent)
				agents[idx].Stats = agentStats
			}
		}
	}

	fmt.Println(ablApps)
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
	fmt.Println("---")
}

func callGetOeManager(instance types.PasInstance, endpoint string) ([]byte, error) {

	client := &http.Client{}
	url := instance.Url + "/oemanager" + endpoint

	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(instance.Creds)))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth) // "Basic YmZ2OmJmdg=="

	res, err := client.Do(req)

	if err == nil {
		body, err := io.ReadAll(res.Body)
		return body, err
	}

	return nil, err
}

func getAblApps(instance types.PasInstance) []types.ABLApplication {

	ablApps := []types.ABLApplication{}

	body, err := callGetOeManager(instance, "/applications")

	if err == nil {
		jsonParsed, _ := gabs.ParseJSON(body)

		for _, child := range jsonParsed.Path("result.Application").Children() {
			var app types.ABLApplication
			json.Unmarshal(child.Bytes(), &app)

			ablApps = append(ablApps, app)
		}
	}

	return ablApps
}

func getAgents(instance types.PasInstance, ablApp types.ABLApplication) []types.MSAgent {

	agents := []types.MSAgent{}

	path := fmt.Sprintf("/applications/%s/agents", ablApp.Name)
	body, err := callGetOeManager(instance, path)

	if err == nil {
		jsonParsed, _ := gabs.ParseJSON(body)

		for _, child := range jsonParsed.Path("result.agents").Children() {

			var agent types.MSAgent
			json.Unmarshal(child.Bytes(), &agent)

			agents = append(agents, agent)
		}
	}

	return agents
}

func getAgentStats(instance types.PasInstance, ablApp types.ABLApplication, agent types.MSAgent) types.MSAgentStat {

	var agentStat types.MSAgentStat

	path := fmt.Sprintf("/applications/%s/agents/%s/metrics", ablApp.Name, agent.AgentId)
	body, err := callGetOeManager(instance, path)

	if err == nil {
		jsonParsed, _ := gabs.ParseJSON(body)
		if len(jsonParsed.Path("result.AgentStatHist").Children()) > 0 {
			jsonBytes := jsonParsed.Path("result.AgentStatHist").Children()[0].Bytes()
			json.Unmarshal(jsonBytes, &agentStat)

			pid, _ := strconv.Atoi(agent.Pid)
			proc, err := process.NewProcess(int32(pid))

			if err == nil {
				mem, _ := proc.MemoryInfo()
				agentStat.MemRSS = mem.RSS
				agentStat.MemVirtual = mem.VMS
			}
		}
	}

	return agentStat
}
