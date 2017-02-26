package metric

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	meminfo       = "/proc/meminfo"
	memOutputFile = "mem"
)

type Memory struct {
	saver
	config                      *Config
	currentMeasure, lastMeasure *memoryMeasure

	// TODO: ajouter DeltaMemFree, DeltaMemOccupied, DeltaSwapFree, ...
}

type memoryMeasure struct {
	MemTotal        kbyte `json:"total"`
	MemFree         kbyte `json:"free"`
	MemOccupied     kbyte `json:"occupied"`
	MemAvailable    kbyte `json:"available"`
	SwapTotal       kbyte `json:"swap-total"`
	SwapFree        kbyte `json:"swap-free"`
	SwapOccupied    kbyte `json:"swap-occupied"`
	VmallocTotal    kbyte `json:"vm-allocated-total"`
	VmallocFree     kbyte `json:"vm-allocated-free"`
	VmallocOccupied kbyte `json:"vm-allocated-occupied"`
}

func NewMemory(config *Config) (*Memory, error) {
	mem := &Memory{}

	saver, err := newSaver(config, mem, memOutputFile)
	if err != nil {
		return nil, err
	}

	mem.saver = *saver
	mem.config = config
	mem.currentMeasure = &memoryMeasure{}
	mem.lastMeasure = &memoryMeasure{}

	return mem, nil
}

func (m *Memory) Update() error {
	*m.lastMeasure = *m.currentMeasure

	return m.currentMeasure.update()
}

func (m *Memory) MarshalCSV() ([]byte, error) {
	return m.currentMeasure.MarshalCSV()
}

func (m *Memory) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.currentMeasure)
}

func (m *Memory) PercentMemFree() float64 {
	return 100.0 - m.PercentMemOccupied()
}

func (m *Memory) PercentMemOccupied() float64 {
	return float64(m.currentMeasure.MemOccupied) * 100.0 /
		float64(m.currentMeasure.MemTotal)
}

func (m *Memory) PercentSwapFree() float64 {
	return 100.0 - m.PercentSwapOccupied()
}

func (m *Memory) PercentSwapOccupied() float64 {
	return float64(m.currentMeasure.SwapOccupied) * 100.0 /
		float64(m.currentMeasure.SwapTotal)
}

func (m *Memory) PercentVmallocFree() float64 {
	return 100.0 - m.PercentVmallocOccupied()
}

func (m *Memory) PercentVmallocOccupied() float64 {
	return float64(m.currentMeasure.VmallocOccupied) * 100.0 /
		float64(m.currentMeasure.VmallocFree)
}

func (m *Memory) String() string {
	format := "\t========== MEMORY ==========\n\n"
	format += "MemTotal:\t %s\n"
	format += "MemFree:\t %s\t%.3f %%\t(%s)\n"
	format += "MemOccupied:\t %s\t%.3f %%\t(%s)\n"
	format += "MemAvailable:\t %s\t\t\t(%s)\n"
	format += "SwapTotal:\t %s\n"
	format += "SwapFree:\t %s\t%.3f %%\t(%s)\n"
	format += "SwapOccupied:\t %s\t%.3f %%\t(%s)\n"
	format += "VmallocTotal:\t %s\n"
	format += "VmallocFree:\t %s\t%.3f %%\t(%s)\n"
	format += "VmallocOccupied: %s\t%.3f %%\t\t(%s)"

	return fmt.Sprintf(format,
		m.currentMeasure.MemTotal,
		m.currentMeasure.MemFree, m.PercentMemFree(), m.currentMeasure.MemFree-m.lastMeasure.MemFree,
		m.currentMeasure.MemOccupied, m.PercentMemOccupied(), m.currentMeasure.MemOccupied-m.lastMeasure.MemOccupied,
		m.currentMeasure.MemAvailable, m.currentMeasure.MemAvailable-m.lastMeasure.MemAvailable,
		m.currentMeasure.SwapTotal,
		m.currentMeasure.SwapFree, m.PercentSwapFree(), m.currentMeasure.SwapFree-m.lastMeasure.SwapFree,
		m.currentMeasure.SwapOccupied, m.PercentSwapOccupied(), m.currentMeasure.SwapOccupied-m.lastMeasure.SwapOccupied,
		m.currentMeasure.VmallocTotal,
		m.currentMeasure.VmallocFree, m.PercentVmallocFree(), m.currentMeasure.VmallocFree-m.lastMeasure.VmallocFree,
		m.currentMeasure.VmallocOccupied, m.PercentVmallocOccupied(), m.currentMeasure.VmallocOccupied-m.lastMeasure.VmallocOccupied)
}

// update updates memoryMeasure parsing /proc/meminfo.
func (m *memoryMeasure) update() error {
	file, err := os.Open(meminfo)
	if err != nil {
		return err
	}
	defer file.Close()

	var n int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "MemTotal") {
			n, err = fmt.Sscanf(line, "MemTotal: %d kB", &m.MemTotal)
			checkSscanf("MemTotal", err, n, 1)
		} else if strings.Contains(line, "MemFree") {
			n, err = fmt.Sscanf(line, "MemFree: %d kB", &m.MemFree)
			checkSscanf("MemFree", err, n, 1)
		} else if strings.Contains(line, "MemAvailable") {
			n, err = fmt.Sscanf(line, "MemAvailable: %d kB", &m.MemAvailable)
			checkSscanf("MemAvailable", err, n, 1)
		} else if strings.Contains(line, "SwapTotal") {
			n, err = fmt.Sscanf(line, "SwapTotal: %d kB", &m.SwapTotal)
			checkSscanf("SwapTotal", err, n, 1)
		} else if strings.Contains(line, "SwapFree") {
			n, err = fmt.Sscanf(line, "SwapFree: %d kB", &m.SwapFree)
			checkSscanf("SwapFree", err, n, 1)
		} else if strings.Contains(line, "VmallocTotal") {
			n, err = fmt.Sscanf(line, "VmallocTotal: %d kB", &m.VmallocTotal)
			checkSscanf("VmallocTotal", err, n, 1)
		} else if strings.Contains(line, "VmallocUsed") {
			n, err = fmt.Sscanf(line, "VmallocUsed: %d kB", &m.VmallocOccupied)
			checkSscanf("VmallocUsed", err, n, 1)
		}
	}

	m.MemOccupied = m.MemTotal - m.MemFree
	m.SwapOccupied = m.SwapTotal - m.SwapFree
	m.VmallocFree = m.VmallocTotal - m.VmallocOccupied

	return nil
}

func (m *memoryMeasure) MarshalCSV() ([]byte, error) {
	str := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d\n",
		m.MemTotal, m.MemFree, m.MemOccupied, m.MemAvailable,
		m.SwapTotal, m.SwapFree, m.SwapOccupied)

	return []byte(str), nil
}

type kbyte int

func (k kbyte) String() string {
	var str string
	fKbyes := float64(k)

	if math.Abs(fKbyes) < 100000 {
		str = fmt.Sprintf("%d kB", int(k))
	} else if math.Abs(fKbyes) < 100000000 {
		str = fmt.Sprintf("%.3f MB", fKbyes/float64(1000))
	} else {
		str = fmt.Sprintf("%.3f GB", fKbyes/float64(1000000))
	}

	return str
}
