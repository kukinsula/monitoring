package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kukinsula/monitoring/metric"
)

type Monitoring struct {
	config  *Config
	metrics []metric.Metric
}

func NewMonitoring(config *Config) (*Monitoring, error) {
	return &Monitoring{
		config: config,
		metrics: []metric.Metric{
			metric.NewCPU(),
			metric.NewMemory(),
			metric.NewNetwork(),
			metric.NewProcesses(),
		},
	}, nil
}

func (m *Monitoring) Start() error {
	clear()

	for {
		for _, m := range m.metrics {
			m.Update()
			m.Save()
			fmt.Println(m, "\n")
		}

		time.Sleep(time.Second)
		clear()
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
