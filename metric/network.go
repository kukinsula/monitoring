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

	// n.update()

	file, err := os.Open(dev)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	var data [nbNetColumns]int64
	var lineName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, ":") {
			fields := strings.Fields(line)
			lineName = fields[0][:len(fields[0])-1]

			for i := 0; i < len(data); i++ {
				data[i], err = strconv.ParseInt(fields[i+1], 10, 0)
				if err != nil {
					log.Fatal(err)
				}
			}

			n.measures[lineName] = &networkInterface{lineName, 0.0, 0.0, data}
			fmt.Println("ici: ", n.measures[lineName].download)

			if n.lastMeasures[lineName] != nil {
				n.measures[lineName].download = n.lastMeasures[lineName].download - float64(data[0]) // float64(1000000)
			}
		}
	}

	fmt.Println(n.measures)
	fmt.Println(n.lastMeasures)
}

func (n *Network) Save() {
}

func (n Network) String() string {
	return ""
}

type networkInterface struct {
	name             string
	download, upload float64
	measure          [nbNetColumns]int64
}

// type networkMeasure struct {
// 	interfaces map[string]networkInterface
// }

// func newNetworkMeasure() *networkMeasure {
// 	return &networkMeasure{
// 		interfaces: make(map[string]networkInterface),
// 	}
// }

// func (n *networkMeasure) update() {
// 	file, err := os.Open(dev)
// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// 	defer file.Close()

// 	var data [nbNetColumns]int64
// 	var lineName string
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()

// 		if strings.Contains(line, ":") {
// 			fields := strings.Fields(line)
// 			lineName = fields[0][:len(fields[0])-1]

// 			for i := 0; i < len(data); i++ {
// 				data[i], err = strconv.ParseInt(fields[i+1], 10, 0)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}

// 			n.interfaces[lineName] = networkInterface{lineName, 0.0, 0.0, data}

// 			// nb, err := fmt.Sscanf(line, "%s: %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d", &lineName, &data[0], &data[1], &data[2], &data[3], &data[4], &data[5], &data[6], &data[7], &data[8], &data[9], &data[10], &data[11], &data[12], &data[13], &data[14], &data[15])
// 		}
// 	}
// }

func isInterface(str string) bool {
	valid := false

	switch str {
	case "wlan0", "l0":
		valid = true
	}

	return valid
}
