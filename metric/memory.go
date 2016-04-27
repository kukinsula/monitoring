package metric

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

var meminfo = "/proc/meminfo"
var memdat = "dat/mem.dat"

func init() {
	_ = os.Remove(memdat)

	file, err := os.Create(memdat)
	if err != nil {
		logger.Fatal(err)
	}

	file.Close()
}

type Memory struct {
	*memoryMeasure
	lastMeasure *memoryMeasure

	// Ajouter DeltaMemFree, DeltaMemOccupied, DeltaSwapFree, ...
}

func NewMemory() *Memory {
	return &Memory{
		memoryMeasure: &memoryMeasure{},
		lastMeasure:   &memoryMeasure{},
	}
}

func (m *Memory) Update() {
	*m.lastMeasure = *m.memoryMeasure
	m.update()
}

func (m Memory) Save() {
	m.save()
}

func (m Memory) PercentMemFree() float64 {
	return 100.0 - m.PercentMemOccupied()
}

func (m Memory) PercentMemOccupied() float64 {
	return float64(m.MemOccupied) * 100.0 / float64(m.MemTotal)
}

func (m Memory) PercentSwapFree() float64 {
	return 100.0 - m.PercentSwapOccupied()
}

func (m Memory) PercentSwapOccupied() float64 {
	return float64(m.SwapOccupied) * 100.0 / float64(m.SwapTotal)
}

func (m Memory) PercentVmallocFree() float64 {
	return 100.0 - m.PercentVmallocOccupied()
}

func (m Memory) PercentVmallocOccupied() float64 {
	return float64(m.VmallocOccupied) * 100.0 / float64(m.VmallocFree)
}

func (m Memory) String() string {
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
		m.MemTotal,
		m.MemFree, m.PercentMemFree(), m.MemFree-m.lastMeasure.MemFree,
		m.MemOccupied, m.PercentMemOccupied(), m.MemOccupied-m.lastMeasure.MemOccupied,
		m.MemAvailable, m.MemAvailable-m.lastMeasure.MemAvailable,
		m.SwapTotal,
		m.SwapFree, m.PercentSwapFree(), m.SwapFree-m.lastMeasure.SwapFree,
		m.SwapOccupied, m.PercentSwapOccupied(), m.SwapOccupied-m.lastMeasure.SwapOccupied,
		m.VmallocTotal,
		m.VmallocFree, m.PercentVmallocFree(), m.VmallocFree-m.lastMeasure.VmallocFree,
		m.VmallocOccupied, m.PercentVmallocOccupied(), m.VmallocOccupied-m.lastMeasure.VmallocOccupied)
}

type memoryMeasure struct {
	MemTotal, MemFree, MemOccupied             kbyte
	SwapTotal, SwapFree, SwapOccupied          kbyte
	VmallocTotal, VmallocFree, VmallocOccupied kbyte
	MemAvailable                               kbyte
}

func (m *memoryMeasure) update() {
	file, err := os.Open(meminfo)
	if err != nil {
		logger.Fatal(err)
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
}

func (m memoryMeasure) save() {
	file, err := os.OpenFile(memdat, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logger.Println(err)
	}
	defer file.Close()

	str := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d\n",
		m.MemTotal, m.MemFree, m.MemOccupied, m.MemAvailable,
		m.SwapTotal, m.SwapFree, m.SwapOccupied)

	w := bufio.NewWriter(file)
	w.WriteString(str)
	w.Flush()
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
