package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kukinsula/monitoring/metric"
)

var (
	supportedMetrics    = []string{"cpu", "mem", "net"}
	nbSupprortedMetrics = len(supportedMetrics)
)

type Monitoring struct {
	config  *metric.Config
	metrics []metric.Metric
}

func NewMonitoring(config *metric.Config) (*Monitoring, error) {
	var fields []string

	// Création du dossier de output
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	config.OutputDir = dir + string(filepath.Separator) +
		config.OutputDir + string(filepath.Separator)

	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return nil, err
	}

	if config.Metrics == "" {
		// TODO : trouver mieux
		fields = supportedMetrics
	} else {
		fields = strings.Split(config.Metrics, ",")

		if len(fields) > nbSupprortedMetrics {
			return nil, fmt.Errorf("too much metrics: max is %d", nbSupprortedMetrics)
		}
	}

	metrics := make([]metric.Metric, 0, nbSupprortedMetrics)

	var m metric.Metric

	for _, field := range fields {
		switch field {
		case "cpu":
			m, err = metric.NewCPU(config)
		case "mem":
			m, err = metric.NewMemory(config)
		case "net":
			m, err = metric.NewNetwork(config)
		case "proc":
			// m, err = metric.NewProcesses(config)
		default:
			err = fmt.Errorf("invalid metric '%s'", field)
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
		for _, metric := range m.metrics {
			err = metric.Update()
			if err != nil {
				return fmt.Errorf("metric update failed: %s", err)
			}

			err = metric.Save()
			if err != nil {
				return fmt.Errorf("metric save failed: %s", err)
			}

			fmt.Printf("%s\n", metric)
		}

		time.Sleep(time.Second)
		clear()
	}
}

func (m *Monitoring) Close() {
	for _, metric := range m.metrics {
		metric.Close()
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
