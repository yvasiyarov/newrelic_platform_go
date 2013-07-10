package main;
//package newrelic_platform_go;

import (
    "log"
    "os"
    "encoding/json"
)

type MetricCollector struct {}

type AgentRunner struct {}

func NewAgent(Version string) * Agent {
    agent := &Agent{
        Version: Version,
    }
    return agent
}

type Agent struct {
    Host string `json:"host"`
    Version string `json:"version"`
    Pid  int `json:"pid"`
}

func (agent *Agent) getEnvironmentInfo() {
    var err error
    agent.Pid = os.Getpid();
    if agent.Host, err = os.Hostname(); err != nil {
        log.Fatalf("Can not get hostname: %#v \n", err)
    }
}

type Component struct {
    Name string `json:"name"`
    GUID string `json:"guid"`
    Duration int `json:"duration"`
    Metrics map[string]interface{} `json:"metrics"`
}

type NewrelicPlugin struct {
    Agent *Agent `json:"agent"`
    Components []Component `json:"components"`
}

func NewNewrelicPlugin() *NewrelicPlugin {
    plugin := &NewrelicPlugin{
        Agent: NewAgent(),
        Components: []Component{},
    }
    return plugin
}

func main() {
    agent := NewAgent("0.0.1")
    agent.getEnvironmentInfo()
    res, _ := json.MarshalIndent(agent, "", "    ");
    
    log.Printf(string(res));
}
