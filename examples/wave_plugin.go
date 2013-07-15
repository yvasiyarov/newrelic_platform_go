package main

import (
	"github.com/yvasiyarov/newrelic_platform_go"
)

type WaveMetrica struct {
	sawtoothMax       int
	sawtoothCounter   int
	squarewaveMax     int
	squarewaveCounter int
}

func (metrica WaveMetrica) GetName() string {
	return "Wave_Metrica"
}
func (metrica WaveMetrica) GetUnits() string {
	return "Queries/Second"
}
func (metrica WaveMetrica) GetValue() (float64, error) {
	return 5, nil
}

func main() {
	plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.1", "7bceac019c7dcafae1ef95be3e3a3ff8866de246", 60)
	component := newrelic_platform_go.NewPluginComponent("Wave component", "com.exmaple.plugin.gowave")
	plugin.AddComponent(component)

	m := WaveMetrica{}
	component.AddMetrica(m)

	plugin.Verbose = true
	plugin.Run()
}
