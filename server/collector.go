package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/segmentio/ksuid"
	"github.com/shirou/gopsutil/process"

	"github.com/bfv/pascollector/misc"
	"github.com/bfv/pascollector/types"
)

func CollectData() []types.Metric {

	var ablApps []types.ABLApplication

	metrics := []types.Metric{}

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

		if len(ablApps) > 0 {

			metric := types.Metric{}
			metric.Id = ksuid.New().String()
			metric.TimeStamp = time.Now().Format(time.RFC3339)
			metric.Instance = instance.Name
			metric.Server = Config.Server
			metric.Metrics = ablApps

			//fmt.Println(metric)
			for _, ablApp := range ablApps {
				fmt.Printf("%s: instance: %s, ablApp: %s, agents: %d\n", metric.Id, instance.Name, ablApp.Name, len(ablApp.Agents))
			}

			metrics = append(metrics, metric)

		} else {
			fmt.Println("server not running")
		}

	}

	return metrics
}

func callGetOeManager(instance types.PasInstance, endpoint string) ([]byte, error) {

	client := &http.Client{}
	url := instance.Url + "/oemanager" + endpoint

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", getAuthString(instance)) // "Basic YmZ2OmJmdg=="

	res, err := client.Do(req)

	if err == nil {
		body, err := io.ReadAll(res.Body)
		return body, err
	}

	return nil, err
}

func getAuthString(instance types.PasInstance) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(misc.Decrypt(instance.Creds)))
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
		children := jsonParsed.Path("result.AgentStatHist").Children()

		if len(children) > 0 {
			jsonBytes := children[0].Bytes()
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
