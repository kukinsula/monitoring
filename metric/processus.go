package metric

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var procdir = "/proc"
var statfile = "stat"

type Processes struct {
	Processes processes
}

type processes []Process

type Process struct {
	Pid, Ppid, Pgrp, Nice, NumThreads int
	Name, State                       string
}

func NewProcesses() *Processes {
	return &Processes{
		Processes: make(processes, 10),
	}
}

func (p *Processes) Update() {
	files, err := ioutil.ReadDir(procdir)
	if err != nil {
		logger.Fatal(err)
	}

	var validProcessID = regexp.MustCompile(`^[0-9]*$`)

	for _, file := range files {
		if file.IsDir() && validProcessID.MatchString(file.Name()) {
			p.readStatPid(file.Name())
		}
	}

	sort.Sort(p.Processes)
}

func (p *Processes) readStatPid(pid string) {
	file, err := os.Open(procdir + "/" + pid + "/" + statfile)
	if err != nil {
		return
	}
	defer file.Close()

	var process = Process{}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	fields := strings.Fields(line)

	process.Pid, err = strconv.Atoi(fields[0])
	if err != nil {
		return
	}
	process.Name = fields[1]
	process.State = fields[2]
	process.Ppid, err = strconv.Atoi(fields[3])
	if err != nil {
		return
	}
	process.Pgrp, err = strconv.Atoi(fields[4])
	if err != nil {
		return
	}
	process.Nice, err = strconv.Atoi(fields[20])
	if err != nil {
		return
	}
	process.NumThreads, err = strconv.Atoi(fields[21])
	if err != nil {
		return
	}

	p.Processes = append(p.Processes, process)
}

func (p Processes) Save() {}

func (p Processes) String() string {
	str := "\t========= PORCESS ==========\n"

	for _, v := range p.Processes {
		str += fmt.Sprintf("%+v\n", v)
	}

	return str
}

func (p processes) Len() int           { return len(p) }
func (p processes) Less(i, j int) bool { return p[i].Pid < p[j].Pid }
func (p processes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
