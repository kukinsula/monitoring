package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kukinsula/monitoring/metric"
)

func main() {
	if len(os.Args) > 1 {
		parseArgs()
	}

	monitor()
}

func monitor() {
	mem := metric.NewMemory()
	cpu := metric.NewCPU()

	clear()

	for {
		cpu.Update()
		mem.Update()

		cpu.Save()
		mem.Save()

		fmt.Println(cpu)
		fmt.Println(mem)

		time.Sleep(time.Second)
		clear()
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func parseArgs() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			help()
			os.Exit(0)
		default:
			help()
			os.Exit(1)
		}
	}
}

func help() {
	fmt.Println(" Monitoring - help")
	fmt.Println("* -h | --help : print help")
}
