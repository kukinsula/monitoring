package metric

import (
	_ "encoding/json"
	"fmt"
	_ "io/ioutil"
	"os"
)

const (
	// TODO : l'utiliser
	CSVSeparator = ","
	fileOpenMode = os.O_CREATE | os.O_WRONLY | os.O_APPEND | os.O_TRUNC
)

type Metric interface {
	Update() error
	Save() error
	Close() error
}

type marshaler interface {
	MarshalJSON() ([]byte, error)
	MarshalCSV() ([]byte, error)
}

type Mode string

var (
	ModeCSV  = Mode("csv")
	ModeJSON = Mode("json")
	ModeWEB  = Mode("web")
)

func (m Mode) GetExtension() string {
	return string(m)
}

type saver struct {
	file      *os.File
	mode      Mode
	marshaler marshaler
}

func newSaver(config *Config, marshaler marshaler, fileName string) (*saver, error) {
	fileName =
		config.OutputDir + fileName + "." + config.Mode.GetExtension()

	_ = os.Remove(fileName)
	file, err := os.OpenFile(fileName, fileOpenMode, 0666)
	if err != nil {
		return nil, err
	}

	return &saver{
		file:      file,
		marshaler: marshaler,
		mode:      config.Mode,
	}, nil
}

func (s *saver) Save() (err error) {
	var b []byte

	switch s.mode {
	case ModeCSV:
		b, err = s.marshaler.MarshalCSV()

	case ModeJSON:
		first := true

		info, err := s.file.Stat()
		if err != nil {
			return err
		}

		size := info.Size()
		if err != nil {
			return err
		}

		if size > 0 {
			first = false

			err = s.file.Truncate(size - 1)
			if err != nil {
				return err
			}
		}

		if first {
			b = append(b, byte('['))
		} else {
			b = append(b, byte(','))
		}

		//
		js, err := s.marshaler.MarshalJSON()
		if err != nil {
			return err
		}

		b = append(b, js...)
		b = append(b, byte(']'))
	default:
		// Ne devrait pas arriver puisque vérifier en amont
		err = fmt.Errorf("invalid mode '%s'", string(s.mode))
	}

	if err != nil {
		return err
	}

	_, err = s.file.Write(b)
	return err
}

func (s *saver) Close() error {
	return s.file.Close()
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
