package metric

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var dev = "/proc/net/dev"
var netdat = "dat/net.dat"

const nbNetColumns = 16

type Network struct {
	*networkMeasure
}

func NewNetwork() *Network {
	return &Network{
		networkMeasure: newNetworkMeasure(),
	}
}

func (n *Network) Update() {
	n.update()
}

func (n *Network) Save() {
}

func (n Network) String() string {
	return ""
}

type networkMeasure struct {
	interfaces map[string][nbNetColumns]int
}

func newNetworkMeasure() *networkMeasure {
	return &networkMeasure{
		interfaces: make(map[string][nbNetColumns]int),
	}
}

func (n *networkMeasure) update() {
	file, err := os.Open(dev)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data [nbNetColumns]int
	var lineName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, ":") {
			fmt.Println(line)

			nb, err := fmt.Sscanf(line, "%s: %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d", &lineName,
				data[0], data[1], data[2], data[3],
				data[4], data[5], data[6], data[7],
				data[8], data[9], data[10], data[11],
				data[12], data[13], data[14], data[15],
			)

			// %d %d %d %d %d %d %d %d %ds %d %d %d %d %d %d %d", _, _, _, _, _, _, _, _, _, _, _, _, _, _, _)

			checkSscanf(lineName, err, nb, 17)
		}
	}
}

func isInterface(str string) bool {
	valid := false

	switch str {
	case "wlan0", "l0":
		valid = true
	}

	return valid
}
