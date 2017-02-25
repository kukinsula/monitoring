package main

import (
	"fmt"
	"os"

	"github.com/kukinsula/monitoring/metric"
)

func main() {
	config, err := metric.NewConfig()
	if err != nil {
		fmt.Printf("NewConfig failed: %s", err)
		os.Exit(1)
	}

	monitoring, err := NewMonitoring(config)
	if err != nil {
		fmt.Printf("NewMonitoring failed: %s", err)
		os.Exit(1)
	}

	monitoring.Start()
}
