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
