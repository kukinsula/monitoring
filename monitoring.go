package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kukinsula/monitoring/metric"
)

type Monitoring struct {
	config  *metric.Config
	metrics []metric.Metric
}

func NewMonitoring(config *metric.Config) (*Monitoring, error) {
	var fields []string

	if config.Metrics == "" {
		// TODO : trouver mieux
		fields = []string{"cpu", "mem", "net", "proc"}
	} else {
		fields = strings.Split(config.Metrics, ",")
	}

	metrics := make([]metric.Metric, 0, len(fields))

	var m metric.Metric
	var err error

	for _, field := range fields {
		switch field {
		case "cpu":
			m, err = metric.NewCPU(config)
		case "mem":
			m, err = metric.NewMemory(config)
		case "net":
			m, err = metric.NewNetwork(config)
		case "proc":
			m, err = metric.NewProcesses(config)
		}

		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	return &Monitoring{config, metrics}, nil
}

func (m *Monitoring) Start() (err error) {
	clear()

	for {
		for _, m := range m.metrics {
			err = m.Update()
			if err != nil {
				return fmt.Errorf("Metric update failed: %s", err)
			}

			err = m.Save()
			if err != nil {
				return fmt.Errorf("Metric save failed: %s", err)
			}

			fmt.Printf("%s\n", m)
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
