package main;
//package newrelic_platform_go;

import (
    "log"
    "os"
    "encoding/json"
    "net/http"
    "time"
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

type MetricaValue interface{}
type SimpleMetricaValue float64
type AggregatedMetricaValue struct {
    Min float64          `json:"min"`
    Max float64          `json:"max"`
    Total float64        `json:"total"`
    Count int            `json:"count"`
    SumOfSquares float64 `json:"sum_of_squares"`
}
func NewAggregatedMetricaValue(existValue float64, newValue float64) *AggregatedMetricaValue {
    v = &AggregatedMetricaValue{
        Min: math.Min(newValue, existVal),
        Max: math.Max(newValue, existVal),
        Total: newValue + existVal,
        Count: 2,
        SumOfSquares: newValue * newValue + existValue * existValue
    }
    return v
}

func (aggregatedValue *AggregatedMetricaValue) Aggregate(newValue float64) {
    aggregatedValue.Min = math.Min(newValue, aggregatedValue.Min)
    aggregatedValue.Max = math.Max(newValue, aggregatedValue.Max)
    aggregatedValue.Total += newValue
    aggregatedValue.Count++
    aggregatedValue.SumOfSquares += newValue * newValue
}

type Metrica interface{
    GetValue() (float64, error)
    GetName() string
}

type Component struct {
    Name string `json:"name"`
    GUID string `json:"guid"`
    Duration int `json:"duration"`
    Metrics map[string]MetricaValue `json:"metrics"`
}

type NewrelicPlugin struct {
    Agent *Agent `json:"agent"`
    Components []Component `json:"components"`
    MetricaModels []Metrica `json:"-"`
    LastPollTime time.Time `json:"-"`
}

func NewNewrelicPlugin() *NewrelicPlugin {
    plugin := &NewrelicPlugin{}

    plugin.Agent = NewAgent(plugin.GetVersion())
    plugin.Agent.CollectEnvironmentInfo()

    component := Component{
        GUID: plugin.GetGuid(),
        Name: plugin.GetPluginName(),
        Duration: plugin.GetReportIntervalInSeconds(),
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

func (plugin *NewrelicPlugin) Harvest() error {
    startTime := time.Now()
    
    if plugin.LastPollTime == 0 {
        plugin.Components[0].Duration = plugin.GetReportIntervalInSeconds()
    } else {
        plugin.Components[0].Duration = startTime.Sub(plugin.LastPollTime)
    }

    plugin.Components[0].Metrics = make(map[string]MetricaValue, len(plugin.MetricaModels))
    for i := 0; i < len(plugin.MetricaModels); i++ {
        model := plugin.MetricaModels[i]
        if newValue, err := model.GetValue(); err == nil {
            if existMetric, ok := plugin.Components[0].Metrics[model.GetName()], ok {
                if floatExistVal, ok := existMetric.(float64); ok {
                    plugin.Components[0].Metrics[model.GetName()] = NewAggregatedMetricaValue(floatExistVal, newValue)
                }
            } else {
                plugin.Components[0].Metrics[model.GetName()].Aggregate(newValue)
            }
        } else {
            log.Printf("Can not get metrica: %v, got error:%#v", model.GetName(), err)
        }
    }

    if httpCode, err := plugin.SendMetricas(); err != nil {
        log.Printf("Can not send metricas to newrelic: %#v\n", err)
    } else {
        if err, isFatal := plugin.CheckResponse(httpCode); isFatal {
            log.Fatalf("Got fatal error:%v\n")
        } else {
            log.Printf("WARNING: %v", err)
        }
    }
    return err
}

func (plugin *NewrelicPlugin) SendMetricas() int, error {
    res, _ := json.MarshalIndent(plugin, "", "    ");
    log.Printf(string(res));
    return 200, nil
}

func (plugin *NewrelicPlugin) CheckResponse(httpResponseCode int) error, bool {
    isFatal := false
    var err error
    switch httpResponseCode {
        http.StatusOK:{
            plugin.Components[0].Metrics = nil
            plugin.LastPollTime = time.Now()
        }
        http.StatusForbidden:{
            err = fmt.Errorf("Authentication error (no license key header, or invalid license key).\n")
            isFatal = true
        }
        http.StatusBadRequest:{
            err = fmt.Errorf("The request or headers are in the wrong format or the URL is incorrect.\n")
            isFatal = true
        }
        http.StatusNotFound:{
            err = fmt.Errorf("Invalid URL\n")
            isFatal = true
        }
        http.StatusRequestEntityTooLarge:{
            err = fmt.Errorf("Too many metrics were sent in one request, or too many components (instances) were specified in one request, or other single-request limits were reached.\n")
            //discard metrics
            plugin.Components[0].Metrics = nil
            plugin.LastPollTime = time.Now()
        }
        http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout: {
            err = fmt.Errorf("Got %v response code.Metricas will be aggregated")
        }
    }
    return err, isFatal
}

func (plugin *NewrelicPlugin) AddMetrica(model Metrica) {
    plugin.MetricaModels = append(plugin.MetricaModels, model)
}

type WaveMetrica struct {
    sawtoothMax int
    sawtoothCounter int
    squarewaveMax int
    squarewaveCounter int 
}
func (metrica WaveMetrica) GetName() string {
    return "Wave Metrica"
}
func (metrica WaveMetrica) GetValue() (float64, error) {
    return 5,  nil
}

func main() {
    plugin := NewNewrelicPlugin()
    m := WaveMetrica{}
    plugin.AddMetrica(m);
    plugin.Harvest()
}
