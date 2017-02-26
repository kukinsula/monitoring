package metric

import (
	"flag"
	"fmt"
)

var (
	DefaultDuration  = 0 // Infini
	DefaultSleep     = 1
	DefaultMetric    = "" // Tous par défaut
	DefaultMode      = ModeCSV
	DefaultOutputDir = "data/"
	DefaultWebServer = ":8080"

	DefaultConfig = &Config{
		Duration:  DefaultDuration,
		Sleep:     DefaultSleep,
		Metrics:   DefaultMetric,
		Mode:      DefaultMode,
		OutputDir: DefaultOutputDir,
		WebServer: DefaultWebServer,
	}
)

type Config struct {
	Duration  int
	Sleep     int
	Metrics   string
	ModeStr   string // Mode : CSV, JSON, ...
	Mode      Mode
	OutputDir string
	WebServer string
}

func NewConfig() (*Config, error) {
	config := DefaultConfig

	flag.IntVar(&config.Duration, "duration", DefaultDuration,
		"Monitoring duration in seconds (0 is infinite)")
	flag.IntVar(&config.Sleep, "sleep", DefaultSleep,
		"Update frequency in seconds")
	flag.StringVar(&config.Metrics, "metrics", DefaultMetric,
		"Metrics to monitor: cpu,mem,proc,net (comma separated)")
	flag.StringVar(&config.ModeStr, "mode", string(DefaultMode),
		"Output mode: CSV, WEB")
	flag.StringVar(&config.OutputDir, "out-dir", DefaultOutputDir,
		"Output files path")
	flag.StringVar(&config.WebServer, "address", DefaultWebServer,
		"Web server address")

	flag.Parse()

	// Mode de output
	switch config.ModeStr {
	case "csv", "CSV", "Csv":
		config.Mode = ModeCSV
	case "json", "JSON", "js", "JS":
		config.Mode = ModeJSON
	case "web", "WEB", "Web":
		config.Mode = ModeWEB
	default:
		return nil, fmt.Errorf("invalid mode '%s'", config.ModeStr)
	}

	return config, nil
}
