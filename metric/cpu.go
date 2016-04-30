package metric

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

var stat = "/proc/stat"
var cpudat = "dat/cpu.dat"

const nbCpuColumns = 10

func init() {
	_ = os.Remove(cpudat)

	file, err := os.Create(cpudat)
	if err != nil {
		logger.Fatal(err)
	}

	file.Close()
}

type CPU struct {
	*cpuMeasure
	lastMeasure  *cpuMeasure
	LoadAverage  float64
	LoadAverages []float64
	NumCPU       int
}

func NewCPU() *CPU {
	NumCPU := runtime.NumCPU()

	return &CPU{
		NumCPU:       NumCPU,
		cpuMeasure:   newCpuMeasure(),
		lastMeasure:  newCpuMeasure(),
		LoadAverages: make([]float64, NumCPU),
	}
}

func (c *CPU) Update() {
	*c.lastMeasure = *c.cpuMeasure
	copy((*c.lastMeasure).cpus, (*c.cpuMeasure).cpus)
	c.cpuMeasure.cpus = make([][nbCpuColumns]int, runtime.NumCPU()+1)

	c.update()
	c.computeCpuAverages()
}

func (c *CPU) computeCpuAverages() {
	c.LoadAverage = c.computeCpuLoad(c.cpus[0], c.lastMeasure.cpus[0])

	for i := 0; i < runtime.NumCPU(); i++ {
		c.LoadAverages[i] = c.computeCpuLoad(c.cpus[i+1], c.lastMeasure.cpus[i+1])
	}
}

func (c *CPU) computeCpuLoad(first, second [nbCpuColumns]int) float64 {
	numerator := float64((second[0] + second[1] + second[2]) -
		(first[0] + first[1] + first[2]))
	denominator := float64((second[0] + second[1] + second[2] + second[3]) -
		(first[0] + first[1] + first[2] + first[3]))

	return math.Abs(numerator / denominator * 100.0)
}

func (c *CPU) Save() {
	file, err := os.OpenFile(cpudat, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logger.Println(err)
	}
	defer file.Close()

	str := fmt.Sprintf("%.2f,", c.LoadAverage)

	for i := 0; i < c.NumCPU; i++ {
		str += fmt.Sprintf("%.2f", c.LoadAverages[i])

		if i != c.NumCPU-1 {
			str += ","
		}
	}

	str += "\n"

	w := bufio.NewWriter(file)
	w.WriteString(str)
	w.Flush()
}

func (c *CPU) String() string {
	str := "\t========== CPU ==========\n\n"
	str += fmt.Sprintf("CPU: \t\t%.2f %%\n", c.LoadAverage)

	for i := 0; i < c.NumCPU; i++ {
		str += fmt.Sprintf("CPU%d: \t\t%.2f %%\n", i, c.LoadAverages[i])
	}

	str += fmt.Sprintf("\nCtxt: \t\t%d (%d)\n", c.Ctxt, c.Ctxt-c.lastMeasure.Ctxt)
	str += fmt.Sprintf("BootTime: \t%d (%v)\n", c.BootTime, time.Unix(c.BootTime, 0))
	str += fmt.Sprintf("Processes: \t%d\n", c.Processes)
	str += fmt.Sprintf("ProcsBlocked: \t%d\n", c.ProcsBlocked)
	str += fmt.Sprintf("ProcsRunning: \t%d", c.ProcsRunning)

	return str
}

type cpuMeasure struct {
	NumberCpus   int
	cpus         [][nbCpuColumns]int
	Ctxt         int
	BootTime     int64
	Processes    int
	ProcsRunning int
	ProcsBlocked int
}

func newCpuMeasure() *cpuMeasure {
	nbCpu := runtime.NumCPU()
	cpus := make([][nbCpuColumns]int, nbCpu+1)

	return &cpuMeasure{
		NumberCpus: nbCpu,
		cpus:       cpus,
	}
}

func (c *cpuMeasure) update() {
	file, err := os.Open(stat)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	var lineName string
	var n, cpuCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "cpu") {
			n, err = fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d", &lineName,
				&c.cpus[cpuCount][0], &c.cpus[cpuCount][1],
				&c.cpus[cpuCount][2], &c.cpus[cpuCount][3],
				&c.cpus[cpuCount][4], &c.cpus[cpuCount][5],
				&c.cpus[cpuCount][6], &c.cpus[cpuCount][7],
				&c.cpus[cpuCount][8], &c.cpus[cpuCount][9],
			)
			checkSscanf(lineName, err, n, 11)
			cpuCount++
		} else if strings.Contains(line, "ctxt") {
			n, err = fmt.Sscanf(line, "ctxt %d", &c.Ctxt)
			checkSscanf("ctxt", err, n, 1)
		} else if strings.Contains(line, "btime") {
			n, err = fmt.Sscanf(line, "btime %d", &c.BootTime)
			checkSscanf("ctxt", err, n, 1)
		} else if strings.Contains(line, "processes") {
			n, err = fmt.Sscanf(line, "processes %d", &c.Processes)
			checkSscanf("ctxt", err, n, 1)
		} else if strings.Contains(line, "procs_running") {
			n, err = fmt.Sscanf(line, "procs_running %d", &c.ProcsRunning)
			checkSscanf("ctxt", err, n, 1)
		} else if strings.Contains(line, "procs_blocked") {
			n, err = fmt.Sscanf(line, "procs_blocked %d", &c.ProcsBlocked)
			checkSscanf("ctxt", err, n, 1)
		}
	}
}

func (c cpuMeasure) String() string {
	var str string

	for i := 0; i < c.NumberCpus; i++ {
		str += fmt.Sprintf("CPU%d: %v\n", i, c.cpus[i])
	}

	str += fmt.Sprintf("Ctxt: %d\n", c.Ctxt)

	return str
}
