package metric

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Lshortfile)

type Updater interface {
	Update() error
}

type Saver interface {
	Save() error
}

type Metric interface {
	Updater
	Saver
}

func checkSscanf(field string, err error, n, expected int) error {
	if err != nil {
		return fmt.Errorf("Sscanf '%s' failed: %s", field, err)
	}

	if n != expected {
		return fmt.Errorf("Sscanf '%s' parsed %d item(s) but expected %d",
			field, n, expected)
	}

	return nil
}
