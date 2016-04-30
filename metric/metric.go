package metric

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Lshortfile)

type Updater interface {
	Update()
}

type Saver interface {
	Save()
}

type Metric interface {
	Updater
	Saver
}

func checkSscanf(field string, err error, n, expected int) {
	if err != nil {
		logger.Fatalf("Sscanf ", field, ": ", err)
	}

	if n != expected {
		logger.Fatalf("Sscanf '%s' parsed %d item(s) but expected %d",
			field, n, expected)
	}
}
