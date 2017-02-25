package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Duration  int
	Sleep     int
	Mode      string // Mode : CSV, JSON, ...
	OutputDir string
	WebServer string
}

func NewConfig() (*Config, error) {
	config := &Config{}

	flag.IntVar(&config.Duration, "duration", 0,
		"Monitoring duration in seconds (0 is infinite)")
	flag.IntVar(&config.Sleep, "sleep", 1,
		"Update frequency in seconds")
	flag.StringVar(&config.Mode, "mode", "CSV",
		"Output mode: CSV, JSON, HTML")
	flag.StringVar(&config.OutputDir, "out-dir", ".",
		"Output files path")
	flag.StringVar(&config.WebServer, "address", ":8080",
		"Web server address")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	return config, nil
}
