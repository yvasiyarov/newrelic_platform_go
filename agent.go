package main;
//package newrelic_platform_go;

import (
    "log"
    "os"
    "encoding/json"
)

var (
    GUID = "com.example.golang_plugin"
    VERSION = "0.0.1"
    PLUGIN_NAME = "Example Go Plugin"
    REPORT_INTERVAL_IN_SECONDS = 60
)
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

func (agent *Agent) CollectEnvironmentInfo() {
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
    plugin := &NewrelicPlugin{}

    plugin.Agent = NewAgent(plugin.GetVersion())
    plugin.Agent.CollectEnvironmentInfo()

    metrics := make(map[string]interface{}, 0)
    component := Component{
        GUID: plugin.GetGuid(),
        Name: plugin.GetPluginName(),
        Duration: plugin.GetReportIntervalInSeconds(),
        Metrics: metrics,
    }
    plugin.Components = []Component{component} 

    return plugin
}

func (plugin *NewrelicPlugin) GetGuid() string {
    return GUID
}

func (plugin *NewrelicPlugin) GetReportIntervalInSeconds() int {
    return REPORT_INTERVAL_IN_SECONDS
}

func (plugin *NewrelicPlugin) GetVersion() string {
    return VERSION
}

func (plugin *NewrelicPlugin) GetPluginName() string {
    return PLUGIN_NAME
}

func main() {
    plugin := NewNewrelicPlugin()
    res, _ := json.MarshalIndent(plugin, "", "    ");
    
    log.Printf(string(res));
}
