package metric

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	dev           = "/proc/net/dev"
	netOutputFile = "net"
	nbNetColumns  = 16
)

type Network struct {
	saver
	config       *Config
	measures     map[string]*networkInterface
	lastMeasures map[string]*networkInterface
}

type networkInterface struct {
	Name     string              `json:"-"`
	Download float64             `json:"download"`
	Upload   float64             `json:"updaload"`
	Measure  [nbNetColumns]int64 `json:"-"`
}

func NewNetwork(config *Config) (*Network, error) {
	net := &Network{}

	saver, err := newSaver(config, net, netOutputFile)
	if err != nil {
		return nil, err
	}

	net.saver = *saver
	net.config = config
	net.measures = make(map[string]*networkInterface)
	net.lastMeasures = make(map[string]*networkInterface)

	return net, nil
}

func (n *Network) Update() error {
	n.lastMeasures = n.measures
	n.measures = make(map[string]*networkInterface)

	file, err := os.Open(dev)
	if err != nil {
		return err
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

			n.measures[interfaceName] = &networkInterface{
				interfaceName, 0.0, 0.0, data,
			}
		}
	}

	n.computeNetworkSpeed()

	return nil
}

func (n *Network) computeNetworkSpeed() {
	for k, _ := range n.measures {
		if n.lastMeasures[k] != nil {
			n.measures[k].Download = float64(n.measures[k].Measure[0]-n.lastMeasures[k].Measure[0]) / float64(1000000)
			n.measures[k].Upload = float64(n.measures[k].Measure[9]-n.lastMeasures[k].Measure[9]) / float64(1000000)
		}
	}
}

func (n *Network) MarshalCSV() ([]byte, error) {
	str := ""
	i, length := 0, len(n.measures)

	for _, v := range n.measures {
		str += fmt.Sprintf("%f,%f", v.Download, v.Upload)

		if i != length-1 {
			str += ","
		}

		i++
	}
	str += "\n"

	return []byte(str), nil
}

func (n *Network) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.measures)
}

func (n *Network) String() string {
	str := "\t========== NETWORK ==========\n\n"
	for _, v := range n.measures {
		str += fmt.Sprintf("%s:\tDownload: %f MB/s,\tUpload: %f MB/s\n",
			v.Name, v.Download, v.Upload)
	}
	return str
}

func isInterface(str string) bool {
	valid := false

	switch str {
	case "wlan0", "l0":
		valid = true
	}

	return valid
}
