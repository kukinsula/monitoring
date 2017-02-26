package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kukinsula/monitoring/metric"
)

func main() {
	flag.Usage = func() { usage(nil) }

	config, err := metric.NewConfig()
	if err != nil {
		usage(err)
	}

	monitoring, err := NewMonitoring(config)
	if err != nil {
		usage(err)
	}

	err = monitoring.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}
}

func usage(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]\n\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}
