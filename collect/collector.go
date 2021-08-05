package collect

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/shirou/gopsutil/process"

	"github.com/bfv/pascollector/types"
)

func CollectData(config types.ConfigFile) {

	var ablApps []types.ABLApplication
	//fmt.Println("collect: " + time.Now().String())

	for _, instance := range config.PasInstances {
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

func SendData(config types.ConfigFile) {
	fmt.Println("send: " + time.Now().String())
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
		children := jsonParsed.Path("result.agents").Children()
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
