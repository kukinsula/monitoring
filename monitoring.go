package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kukinsula/monitoring/metric"
)

func main() {
	_, err := NewConfig()
	if err != nil {
		fmt.Printf("NewConfig failed: %s", err)
		os.Exit(1)
	}

	monitor()
}

func monitor() {
	var metrics = []metric.Metric{
		metric.NewCPU(),
		metric.NewMemory(),
		metric.NewNetwork(),
		metric.NewProcesses(),
	}

	clear()

	for {
		for _, m := range metrics {
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
