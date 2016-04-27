package metric

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var dev = "/proc/net/dev"
var netdat = "dat/net.dat"

const nbNetColumns = 16

type Network struct {
	measures     map[string]*networkInterface
	lastMeasures map[string]*networkInterface
}

func NewNetwork() *Network {
	return &Network{
		measures:     make(map[string]*networkInterface),
		lastMeasures: make(map[string]*networkInterface),
	}
}

func (n *Network) Update() {
	n.lastMeasures = n.measures
	n.measures = make(map[string]*networkInterface)

	file, err := os.Open(dev)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	var data [nbNetColumns]int64
	var interfaceName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, ":") {
			fields := strings.Fields(line)
			interfaceName = fields[0][:len(fields[0])-1]

			for i := 0; i < len(data); i++ {
				data[i], err = strconv.ParseInt(fields[i+1], 10, 0)
				if err != nil {
					log.Fatal(err)
				}
			}

			n.measures[interfaceName] = &networkInterface{interfaceName, 0.0, 0.0, data}
		}
	}

	n.computeNetworkSpeed()
}

func (n *Network) computeNetworkSpeed() {
	for k, _ := range n.measures {
		if n.lastMeasures[k] != nil {
			n.measures[k].download = float64(n.measures[k].measure[0]-n.lastMeasures[k].measure[0]) / float64(1000000)
			n.measures[k].upload = float64(n.measures[k].measure[9]-n.lastMeasures[k].measure[9]) / float64(1000000)
		}
	}
}

func (n *Network) Save() {
}

func (n Network) String() string {
	str := "\t========== NETWORK ==========\n\n"
	for _, v := range n.measures {
		str += fmt.Sprintf("%s:\tDownload: %f MB/s,\tUpload: %f MB/s\n",
			v.name, v.download, v.upload)
	}
	return str
}

type networkInterface struct {
	name             string
	download, upload float64
	measure          [nbNetColumns]int64
}

func isInterface(str string) bool {
	valid := false

	switch str {
	case "wlan0", "l0":
		valid = true
	}

	return valid
}
