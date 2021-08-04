package types

type PasResponse struct {
	Operation    string                 `json:"operation"`
	Outcome      string                 `json:"outcome"`
	ErrorMessage string                 `json:"errmsg"`
	Version      string                 `json:"versionStr"`
	VersionNo    int                    `json:"versionNo"`
	Result       map[string]interface{} `json:"result"`
}

type ABLApplication struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Agents  []MSAgent
}

type MSAgent struct {
	AgentId string `json:"agentId"`
	Pid     string `json:"pid"`
	State   string `json:"state"`
	Stats   MSAgentStat
}

type MSAgentStat struct {
	ActiveThreads     int    `json:"ActiveThreads"`
	ActiveSessions    int    `json:"ActiveSessions"`
	OpenConnections   int    `json:"OpenConnections"`
	ExitedThreads     int    `json:"ExitedThreads"`
	ExitedSessions    int    `json:"ExitedSessions"`
	ClosedConnections int    `json:"ClosedConnections"`
	CStackMemory      uint64 `json:"CStackMemory"`
	OverheadMemory    uint64 `json:"OverheadMemory"`
	MemRSS            uint64
	MemVirtual        uint64
}
