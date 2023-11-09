package simpleconfig

import (
	"bytes"
	"os"
	"tsundere/packages/utilities/simpleconfig/encoders"
)

func (s *SimpleConfig) decode(path string, v interface{}) error {
	err := s.encode(path, false, v)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = encoders.NewDecoder(int(s.coder), bytes.NewReader(file)).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
